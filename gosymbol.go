package gosymbol

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
)

/* Evaluation */

func (e undefined) Eval() Func {
	return func(args Arguments) float64 {return math.NaN()}
}

func (e constant) Eval() Func {
	return func(args Arguments) float64 {return e.Value}
}

func (e variable) Eval() Func {
	return func(args Arguments) float64 {return args[e.Name]}
}

func (e add) Eval() Func {
	return func(args Arguments) float64 {
		sum := e.Operands[0].Eval()(args) // Initiate with first operand since 0 may not always be identity
		for ix := 1; ix < len(e.Operands); ix++ {
			sum += e.Operands[ix].Eval()(args)
		}
		return sum
	}
}

func (e mul) Eval() Func {
	return func(args Arguments) float64 {
		prod := e.Operands[0].Eval()(args) // Initiate with first operand since 1 may not always be identity
		for ix := 1; ix < len(e.Operands); ix++ {
			prod *= e.Operands[ix].Eval()(args)
		}
		return prod
	}
}

func (e exp) Eval() Func {
	return func(args Arguments) float64 {return math.Exp(e.Arg.Eval()(args))}
}

func (e log) Eval() Func {
	return func(args Arguments) float64 {return math.Log(e.Arg.Eval()(args))}
}

func (e pow) Eval() Func {
	return func(args Arguments) float64 {return math.Pow(e.Base.Eval()(args), e.Exponent.Eval()(args))}
}

/* Evaluation to string */

func (e undefined) String() string {
	return "Undefined"
}

func (e constant) String() string {
	if e.Value < 0 {
		return fmt.Sprintf("( %v )", e.Value)	
	} else {
		return fmt.Sprint(e.Value)
	}
}

func (e variable) String() string {
	return string(e.Name)
}

func (e add) String() string {
	str := fmt.Sprintf("( %v", e.Operands[0])
	for ix := 1; ix < len(e.Operands); ix++ {
		str += fmt.Sprintf(" + %v", e.Operands[ix])
	}
	str += " )"
	return str
}

func (e mul) String() string {
	str := fmt.Sprintf("( %v", e.Operands[0])
	for ix := 1; ix < len(e.Operands); ix++ {
		str += fmt.Sprintf(" * %v", e.Operands[ix])
	}
	str += " )"
	return str
}

func (e exp) String() string {
	return fmt.Sprintf("exp( %v )", e.Arg)
}

func (e log) String() string {
	return fmt.Sprintf("log( %v )", e.Arg)
}

func (e pow) String() string {
	return fmt.Sprintf("( %v^%v )", e.Base, e.Exponent)
}

/* Helper Functionality */

// Match case function. Performs the specified action when the top operand of `expr` 
// macthes a case.
//
// Note: this function ought to be used for in every place which would require modificaton
// uppon adding new operands / functions to the package. Reason being that it will force 
// a compilation error if you have not added a case handling the new operator everywhere 
// it is needed, instead of having unexpected behaviour during runtime.
func MatchTransform[T any] (
	expr Expr,
	undefAction func(undefined) T,
	constAction func(constant) T,
	varAction func(variable) T,
	addAction func(add) T,
	mulAction func(mul) T,
	powAction func(pow) T,
	expAction func(exp) T,
	logAction func(log) T,
	sqrtAction func(sqrt) T,

) T {
	switch expr := expr.(type) {
	case undefined:
		return undefAction(expr)
	case constant:
		return constAction(expr)
	case variable:
		return varAction(expr)
	case add:
		return addAction(expr)
	case mul:
		return mulAction(expr)
	case pow:
		return powAction(expr)
	case exp:
		return expAction(expr)
	case log:
		return logAction(expr)
	case sqrt:
		return sqrtAction(expr)
	default:
		errMsg := fmt.Sprintf("ERROR: argument of type %v has no implemented match case", reflect.TypeOf(expr))
		panic(errMsg)
	}
}

// Substitutes u for t in expr.
func Substitute(expr, u, t Expr) Expr {
	if Equal(u, t) {
		return u
	} else if Equal(expr, u) {
		return t
	} else if Contains(expr, u) {	
		for ix := 1; ix <= NumberOfOperands(expr); ix++ {
			processedOp := Substitute(Operand(expr, ix), u, t)
			expr = replaceOperand(expr, ix, processedOp)
		}
		// Subsituting here again in order to continue to Substitute
		// if the result from above substitution is something that can be
		// Substituted again. Se test 7 and test 12 in gosymbol_test.go.
		return Substitute(expr, u, t)
	} else {
		return expr
	}
}

/*
Replaces operand number n in t with u and returns the resulting
expression. The function panics if n is larger than 
the NumberOfOperands(t).
*/
func replaceOperand(t Expr, n int, u Expr) Expr {
	nop := NumberOfOperands(t)
	if n > nop {
		errMsg := fmt.Sprintf("ERROR: trying to access operand %v but expr has only %v operands.", n, nop)
		panic(errMsg)
	} else if n <= 0 {
		errMsg := fmt.Sprintf("ERROR: there exists no non-positive indexed operands, you are trying to replace operand: %v", nop)
		panic(errMsg)
	}

	return MatchTransform[Expr](
		t,
		func (expr undefined) Expr {
			return expr
		},
		func(expr constant) Expr {
			return expr
		},
		func(expr variable) Expr {
			return expr
		},
		func(expr add) Expr {
			expr.Operands[n-1] = u
			return expr
		},
		func(expr mul) Expr {
			expr.Operands[n-1] = u
			return expr
		},
		func(expr pow) Expr {
			if n == 1 {
				expr.Base = u
			} else {
				expr.Exponent = u
			} 
			return expr
		},
		func(expr exp) Expr {
			expr.Arg = u
			return expr 
		},
		func(expr log) Expr {
			expr.Arg = u
			return expr
		},
		func(expr sqrt) Expr {
			expr.Arg = u
			return expr
		},
	)

	// Since the first cases has no operands
	// we consider replacing one as just returning
	// the original expression
	//switch v := t.(type) {
	//case undefined:
	//	return v
	//case constant:
	//	return v
	//case variable:
	//	return v
	//case add:
	//	v.Operands[n-1] = u
	//	return v
	//case mul:
	//	v.Operands[n-1] = u
	//	return v
	//case pow:
	//	if n == 1 {
	//		v.Base = u
	//	} else {
	//		v.Exponent = u
	//	}
	//	return v
	//case exp:
	//	v.Arg = u
	//	return v
	//case log:
	//	v.Arg = u
	//	return v
	//case sqrt:
	//	v.Arg = u
	//	return v
	//default:
	//	errMsg := fmt.Sprintf("ERROR: function is not implemented for type: %v", reflect.TypeOf(v))
	//	panic(errMsg)
	//}
}

/*
Swaps operand number n1 with operand number n2 in expr.
*/
func swapOperands(expr Expr, n1, n2 int) Expr {
	op1 := Operand(expr, n1)
	op2 := Operand(expr, n2)
	expr = replaceOperand(expr, n1, op2)
	expr = replaceOperand(expr, n2, op1)
	return expr
}

/* 
Recursively checks exact syntactical equality between t and u,
i.e. it does not simplify any expression nor does it 
take any properties, e.g. commutativity, into account.
*/
func Equal(t, u Expr) bool {
	return MatchTransform[bool](
		t,
		func (expr undefined) bool {
			_, ok := u.(undefined)
			return ok
		},
		func(expr constant) bool {
			uTyped, ok := u.(constant)
			return ok && expr.Value == uTyped.Value
		},
		func(expr variable) bool {
			uTyped, ok := u.(variable)
			return ok && expr.Name == uTyped.Name
		},
		func(expr add) bool {
			_, ok := u.(add)
			if !ok {
				return false
			} else if NumberOfOperands(expr) != NumberOfOperands(u) {
				return false
			}

			for ix := 1; ix <= NumberOfOperands(expr); ix++ {
				exprOp := Operand(expr, ix)
				uOp := Operand(u, ix)
				if !Equal(exprOp, uOp) {
					return false
				}
			}
			return true
		},
		func(expr mul) bool {
			_, ok := u.(mul)
			if !ok {
				return false
			} else if NumberOfOperands(expr) != NumberOfOperands(u) {
				return false
			}

			for ix := 1; ix <= NumberOfOperands(expr); ix++ {
				exprOp := Operand(expr, ix)
				uOp := Operand(u, ix)
				if !Equal(exprOp, uOp) {
					return false
				}
			}
			return true
		},
		func(expr pow) bool {
			_, ok := u.(pow)
			return ok && Equal(Operand(expr, 1), Operand(u, 1)) && Equal(Operand(expr, 2), Operand(u, 2))
		},
		func(expr exp) bool {
			_, ok := u.(exp)
			return ok && Equal(Operand(expr, 1), Operand(u, 1))
		},
		func(expr log) bool {
			_, ok := u.(log)
			return ok && Equal(Operand(expr, 1), Operand(u, 1))
		},
		func(expr sqrt) bool {
			_, ok := u.(sqrt)
			return ok && Equal(Operand(expr, 1), Operand(u, 1))
		},
	)
}

/*
Returns true if t and u are equal up to type 
for every element in resp. syntax tree, otherwise 
false. This means that two constants of different 
value, or variables with different names, would
return true.
*/
func TypeEqual(t, u Expr) bool {
	if !isSameType(t, u) {
		return false
	} else if NumberOfOperands(t) != NumberOfOperands(u) {
		return false
	}

	// Base cases are the leaf node types
	switch t.(type) {
	case undefined:
		_, ok := u.(undefined)
		return ok
	case constant:
		_, ok := u.(constant)
		return ok
	case variable: 
	_, ok := u.(variable)
		return ok
	}

	// Recusively checks if every operand is type equal 
	// as well. Breaks and returns false if any of the
	// operands are not equal.
	// TODO: this does not take associativty of e.g. add and mul into account.
	for ix := 1; ix <= NumberOfOperands(t); ix++ {
		tOperand := Operand(t, ix)
		uOperand := Operand(u, ix)
		if !TypeEqual(tOperand, uOperand) {
			return false
		}
	}

	return true
}

/* 
Recursively checks exact equality between expr and u,
i.e. it does not simplify any expression nor does it 
take any properties, e.g. commutativity, into account.
*/
func Contains(expr, u Expr) bool {
	switch v := expr.(type) {
	case undefined:
		_, ok := u.(undefined)
		return ok
	case constant:
		uTyped, ok := u.(constant)
		return ok && v.Value == uTyped.Value
	case variable:
		uTyped, ok := u.(variable)
		return ok && v.Name == uTyped.Name
	default:
		if Equal(v,u) {return true}
		for ix := 1; ix <= NumberOfOperands(v); ix++ {
			vOp := Operand(v, ix)
			if Contains(vOp, u) {
				return true
			}
		}
		return false
	}
}

/*
Returns the different variable names 
present in the given expression.
*/
func VariableNames(expr Expr) []VarName {
	var stringSlice []string
	variableNames(expr, &stringSlice)
	if len(stringSlice) == 0 {
		return []VarName{}
	}

	// Checks for and deletes duplicates of variable names
	// TODO: speed up by checking for duplicates during sorting
	sort.Strings(stringSlice)
	var variableNamesSlice []VarName
	variableNamesSlice = append(variableNamesSlice, VarName(stringSlice[0]))
	jx := 0
	for ix := 1; ix < len(stringSlice); ix++ {
		if variableNamesSlice[jx] != VarName(stringSlice[ix]) {
			variableNamesSlice = append(variableNamesSlice, VarName(stringSlice[ix]))
			jx++
		}
	}

	return variableNamesSlice
}

/* 
Recursively travesrses the whole AST and appends
the variable names to targetSlice. 
*/
func variableNames(expr Expr, targetSlice *[]string) {
	switch v := expr.(type) {
	case undefined:
		return 
	case constant:
		return
	case variable:
		*targetSlice = append(*targetSlice, string(v.Name))
	default:
		for ix := 1; ix <= NumberOfOperands(expr); ix++ {
			op := Operand(expr, ix)
			variableNames(op, targetSlice)
		}
	}
}

// Returns the number of operands for top level operation.
func NumberOfOperands(expr Expr) int {
	return MatchTransform[int](
		expr, 
		func(expr undefined) int { return 0 }, 
		func(expr constant) int { return 0 }, 
		func(expr variable) int { return 0 }, 
		func(expr add) int { return len(expr.Operands) }, 
		func(expr mul) int { return len(expr.Operands) }, 
		func(expr pow) int { return 2 }, 
		func(expr exp) int { return 1 },
		func(expr log) int { return 1 },
		func(expr sqrt) int { return 1 },
	)
}

/*
Returns the n:th (starting at 1) operand (left to right) of expr.
If expr has no operands it returns nil.
If n is larger than NumberOfOperands(expr) it will panic.
*/
func Operand(expr Expr, n int) Expr {
	nop := NumberOfOperands(expr)
	if n > nop {
		errMsg := fmt.Sprintf("ERROR: trying to access operand %v but expr has only %v operands.", n, nop)
		panic(errMsg)
	}

	return MatchTransform[Expr](
		expr, 
		func(expr undefined) Expr { return nil }, 
		func(expr constant) Expr { return nil }, 
		func(expr variable) Expr { return nil }, 
		func(expr add) Expr { return expr.Operands[n-1] }, 
		func(expr mul) Expr { return expr.Operands[n-1] }, 
		func(expr pow) Expr { 
			if n == 1 {
				return expr.Base
			} else {
				return expr.Exponent
			}
		}, 
		func(expr exp) Expr { return expr.Arg },
		func(expr log) Expr { return expr.Arg },
		func(expr sqrt) Expr { return expr.Arg },
	)
}

// TODO: see Computer Algebra and Symbolic Computation page 10 to understand this shit
func Map(F Expr, u ...Expr) Expr {panic("Not implemented yet")}


func isSameType(a, b any) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}

// Returns the deepest depth of expr
func Depth(expr Expr) int {
	switch expr.(type) {
	case undefined:
		return 0
	case constant:
		return 0
	case variable:
		return 0
	default:
		maxDepth := 0
		for ix := 1; ix <= NumberOfOperands(expr); ix++ {
			op := Operand(expr, ix)
			opDepth := Depth(op)
			if opDepth > maxDepth {
				maxDepth = opDepth
			}
		}
		return maxDepth + 1
	}
}

/* Automatic Simplification */


/*
Flattens and sorts the top level of a tree sum or product, i.e. 
doing the transofmration:
	( (x + y) + z ) + ( (i + j) + (k + l) ) -> (i + j) + (k + l) + (x + y) + z 
*/
func flattenTopLevel(expr Expr) Expr {
	// If expr is add or mul we check if any of its operands 
	// are of the same type. If so we "bring it up one level"
	// otherwise we do nothing.
	switch expr.(type) {
	case add:
		newTerms := make([]Expr, 0)
		for ix := 1; ix <= NumberOfOperands(expr); ix++ {
			op := Operand(expr, ix)
			if opTyped, ok := op.(add); ok {
				// Here we "bring it up one level" by appending the 
				// operand's operands to the new set of operands
				newTerms = append(newTerms, opTyped.Operands...)
			} else {
				newTerms = append(newTerms, op)	
			}
		}
		return topOperandSort(Add(newTerms...))
	case mul: 
		newFactors := make([]Expr, 0)
		for ix := 1; ix <= NumberOfOperands(expr); ix++ {
			op := Operand(expr, ix)
			if opTyped, ok := op.(mul); ok {
				// Here we "bring it up one level" by appending the 
				// operand's operands to the new set of operands
				newFactors = append(newFactors, opTyped.Operands...)
			} else {
				newFactors = append(newFactors, op)
			}
		}
		return topOperandSort(Mul(newFactors...))
	default:
		return expr
	}
}

/*
Applies rule to expr and returns the transformed expression.
If expression does not match rule the ingoing expression 
will just be returned.
*/
func (rule transformationRule) apply(expr Expr) (Expr, bool) {
	if rule.match(expr) {
		return rule.transform(expr), true
	}
	return expr, false
}

func (rule transformationRule) match(expr Expr) bool {
	// Fisrt check if pattern is defined. If not
	// we execute patternFunction if it exists. 
	// If no pattern or patternFunction exists we return false 
	if rule.pattern != nil {
		varCache := variableCache()
		return patternMatch(rule.pattern, expr, varCache)
	} else if rule.patternFunction != nil {
		return rule.patternFunction(expr)
	} else {
		return false
	}
}

/*
Recursively checks if expr matches pattern. varCache is an empty
map internally used to keep track of what the variables in pattern 
corresponds to in expr. The function expects that no variable has the 
same name as a constrained variable.
*/
func patternMatch(pattern, expr Expr, varCache _variableCache) bool {
	if patternTyped, ok := pattern.(cacheVariable); ok {
			
	}

	return MatchTransform[bool](
		expr, 
		func(expr undefined) bool { // Base case 
			// Trying to match with undefined is alwasy false
			return false 
		},
		func(expr constant) bool { // Base case 
			patternTyped, ok := pattern.(constant)
			return ok && patternTyped.Value == expr.Value
		},
		func(expr variable) bool { // Base case 
			patternTyped, ok := pattern.(variable)
			return ok && patternTyped.Name == expr.Name
		},
		func(expr add) bool { 
			if _, ok := pattern.(add); !ok { return false }
			return patternMatchOperands(pattern, expr, varCache) 
		},
		func(expr mul) bool { 
			if _, ok := pattern.(mul); !ok { return false }
			return patternMatchOperands(pattern, expr, varCache) 
		},
		func(expr pow) bool { 
			if _, ok := pattern.(pow); !ok { return false }
			return patternMatchOperands(pattern, expr, varCache) 
		},
		func(expr exp) bool { 
			if _, ok := pattern.(exp); !ok { return false }
			return patternMatchOperands(pattern, expr, varCache) 
		},
		func(expr log) bool { 
			if _, ok := pattern.(log); !ok { return false }
			return patternMatchOperands(pattern, expr, varCache) 
		},
		func(expr sqrt) bool { 
			if _, ok := pattern.(sqrt); !ok { return false }
			return patternMatchOperands(pattern, expr, varCache) 
		},
	)

	switch v := pattern.(type) {
	case undefined:
		_, ok := expr.(undefined)
		return ok
	case constant:
		exprTyped, ok := expr.(constant)
		return ok && v.Value == exprTyped.Value 
	case cacheVariable:
		// If the cached variable as an expression assigned to it we pattern 
		// match it with `expr`. If it does not and there is no constraint 
		// on the variable defined in the pattern we cache current expression,
		// otherwise we just check if the constraint is satisfied and then cache.
		if !varCache.isUnassigned(v.Name) {
			cachedExpr := varCache.get(v.Name)
			return patternMatch(cachedExpr, expr, varCache)
		} else if v.Constraint == nil {
			varCache.add(v.Name, expr)	
			return true
		} else if v.Constraint(expr) {
			varCache.add(v.Name, expr)
			return true
		}
		return false
	case variable:
		exprTyped, ok := expr.(variable)
		if !ok {
			return false
		}
		return v.Name == exprTyped.Name 
	case add:
		if _, ok := expr.(add); !ok { return false }
		return patternMatchOperands(v, expr, varCache) 
	case mul:
		if _, ok := expr.(mul); !ok { return false }
		return patternMatchOperands(v, expr, varCache)
	case pow:
		if _, ok := expr.(pow); !ok { return false }
		return patternMatchOperands(v, expr, varCache)
	case exp:
		if _, ok := expr.(exp); !ok { return false }
		return patternMatchOperands(v, expr, varCache)
	case log:
		if _, ok := expr.(log); !ok { return false }
		return patternMatchOperands(v, expr, varCache)
	case sqrt:
		if _, ok := expr.(sqrt); !ok { return false }
		return patternMatchOperands(v, expr, varCache)
	default:
		errMsg := fmt.Errorf("ERROR: expression of type: %v have no matchPattern case implemented", reflect.TypeOf(v))
		panic(errMsg)
	}
}

/*
Checks if the operands of pattern and expr match.
This function does not check if the main operator 
of pattern and expr match.
*/
func patternMatchOperands(pattern, expr Expr, varCache variableCache) bool {
	if NumberOfOperands(pattern) != NumberOfOperands(expr) {
		return false
	}

	// Recursively checks if each operand matches
	for ix := 1; ix <= NumberOfOperands(pattern); ix++ {
		patternOp := Operand(pattern, ix)
		exprOp := Operand(expr, ix)
		if !patternMatch(patternOp, exprOp, varCache) {
			return false
		}
	}
	return true
}

// TODO: figure this out
func Expand(expr Expr) Expr {
	return nil
}

/*
Returns a new expression where the
terms in in s1 has been prepended
to the terms in s2, i.e. the order 
of terms is not changed.
*/
func mergeSums(s1, s2 add) add {
	panic("not implemented yet")
}

/*
Returns a new expression where the
factors in in p1 has been prepended
to the factors in p2, i.e. the order 
of factors is not changed.
*/
func mergeProducts(p1, p2 mul) mul {
	panic("not implemented yet")
}

/*
Sorts the operands of expr in increasing order in accordance
with the order relation defined by compare(e1,e2 Expr). It does
not recursively sort the operands operands etc. This should only 
be applied when the operator is commutative!

To make a complete sort using this it needs to be recursive
called on the operands.

NOTE: Using insertion sort so worst case time complexity is O(n^2).
*/
func topOperandSort(expr Expr) Expr {
	for ix := 1; ix <= NumberOfOperands(expr)-1; ix++ {
		op1 := Operand(expr, ix)
		op2 := Operand(expr, ix+1)
		
		n := ix
		for !op1.compare(op2) {
			expr = swapOperands(expr, n, n+1)
			
			// As long as n > 1 we are not at the
			// first operand and there is no risk of 
			// index out of bounds error. If we are at 
			// the last operand we need to break since 
			// compare will continue to return false if 
			// op1 == op2.
			if n > 1 {
				op1 = Operand(expr, n-1)
				n--
			} else {
				break
			}
		}
	}
	return expr
}

func TopOperandSort(expr Expr) Expr { return topOperandSort(expr) }

/* Variable Cacheing */

func variableCache() _variableCache {
	varCache := _variableCache{}
	varCache.cache = make(map[string]Expr)
	return varCache
}

/*
Converts the key used in the API to the internally used hashmap key.
*/
func toInternalKey(externalKey VarName) string { 
	return string(externalKey) + "_INTERVAL" 
}

/*
Converts the internally used hashmap key to the key used in the API.
*/
func toExternalKey(internalKey string) VarName {
	trimmedKey := strings.TrimSuffix(internalKey, "_INTERVAL")
	return VarName(trimmedKey)
}

/*
Adds `expr` to cache as `key`. Note that if `key` is already in vc it will be 
over-written.
*/
func (vc _variableCache) add(key VarName, expr Expr) {
	vc.cache[toInternalKey(key)] = expr	
}

// TODO
func (vc _variableCache) remove(key string, expr Expr) {}

/*
Returns the exression cached as `key`. Note that this method 
assumes that `key` exists in the underlying map.
*/
func (vc _variableCache) get(key VarName) Expr {
	return vc.cache[toInternalKey(key)]
}

// TODO
func (vc _variableCache) isCached(expr Expr) bool {
	panic("not implemented!")
}

func (vc _variableCache) isUnassigned(key VarName) bool {
	_, ok := vc.cache[toInternalKey(key)]
	return !ok
}

