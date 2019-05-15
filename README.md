#Go coding guidelines

##Suggested static code analysis tools
- Code should pass the standard SCA framework `go vet` 
- Code should not be changed after running `gofmt` (running it fixes formatting issues, so it is advised to run this command before delivering the code) as an alternative go lint might be used as SCA
- No code is delivered which fails `go test ./test/…` command in the repository's root folder
- Code should pass [GoCyclo](https://github.com/fzipp/gocyclo), which calculates cyclomatic complexities of functions in Go source code
- Code should pass [Ineffective Assign](https://github.com/gordonklaus/ineffassign), which detects ineffectual assignments in Go code
- Code should pass [Miss Spell](https://github.com/client9/misspell) which detects English spell errors in the source code

##Recommended code and software architecture guidelines

As the main source of requirements we should use [Effective go instructions](https://golang.org/doc/effective_go.html) and [language specifications](https://golang.org/ref/spec).

Additionally we suggest following guidelines.

###External service calls
When calling network services (databases, message brokers etc) we should take care about a delayed availability as well as connection timeouts. We should have a reconnection/recovery policy when an external service is gone or gets restarted.

When calling external services with go routines it’s recommended to use [Context library](https://golang.org/pkg/context) to control timeouts/deadlines and cleanly release resources.

###Go routines
When starting a go routine, we should always manage it’s lifecycle:

Runtime errors should be forwarded through channels

Acquired resources should be released when the main process exists, go routine should have a possibility for resource releasing by getting signals over channels (also consider using context package for that). 

Consider shared resources management, don’t forget the principle of [sharing by communicating rather than communicate by sharing memory](https://golang.org/doc/effective_go.html#sharing)

###Errors handling
No errors should be ignored. We shouldn’t use blank identifiers for error types e.g.

    res, _ := createResource()
    
All errors should be handled either as critical failures with a required exit or warnings with continuation.
Critical errors should be returned, interrupting the current function flow. Entry points (e.g. main function) should always get critical errors from lower layers in return statements. 

All critical errors should be logged with the call stack output for later debugging. Only the entry points decide if we should exit the execution with a non zero code or continue trying with some repetition policy. 
No functions should interrupt the execution on critical errors, which practically means avoiding following code chunks:

    if err != nil {
        panic(err)
    }

And doing this instead:

    if err != nil {
        return err
    }

When a non-critical error appears, we can do the following:

- Log error and continue the normal execution:


        if err != nil {
            log.PrintLn(err)
        }


- Cast error to some recognisable error type (e.g. `WarningError`) and pass it to the outer caller as a return result. The later will handle this error based on it’s type:


        func NewWarningError(err error) error {
        	return WarningError(err: err)
        }
        if err != nil {
        	wrappedError := NewWarningError(err)
        	return wrappedError
        }
        
        //and then in the entry point
        switch err.(type) {
        	case WarningError:
        		log.Error(err)
        	case CriticalError:
        		log.Panic(err)
        	default:
        		log.Panic(err)
        }

###Logging
Every meaningful code operation should be logged. This gives a great help for debugging in production mode, when the program doesn't behave as expected. We should generally log not only what the program is currently doing but also all i/o contexts. 

The recommended way of logging is using [Log package](https://golang.org/pkg/log) rather than [Fmt package](https://golang.org/pkg/fmt) which is less customizable. 

The lower level functions should never change logging format nor logging behaviour (e.g. by calling `SetOutput`).

The default behaviour is to forward info/warning logs to stdout and error logs to stderr. 
The decision of how to handle logging data should be made on the infrastructure level (e.g. classical way to fetch stdouts by Docker and forward them to ElasticSearch with [FileBeat](https://www.elastic.co/products/beats/filebeat). 

###Imports
Blank import identifiers should only be used for external libraries when they require this (e.g. for init function execution).

Dot (.) import notation should not be used.

###External libraries
[Go modules](https://github.com/golang/go/wiki/Modules) is recommended as dependency management tool for all supported go repositories as an obvious and modern way of fetching external dependencies.

Using `dep` and `go get` syntax should be considered as outdated.

###Comments
The general rule is to use no commented out code in production. There are no case where we should keep some “code knowledge” for later use. 

We always can refer to the git history to look at previous code versions.

Inline comments are generally an exception rather than a normal case. The code itself should be written as a documentation (see the motivation here). 

###Names
Package names as brief, lowercase, singular words reflecting the domain which is handled by the code in it. 
There is no need for camel/kebab/underscore concatenated word groups.

Exported structs, functions, variables should not contain words used in package names (e.g. `bufio.Reader` rather than `bufio.BufReader`)

Use camelcase for function names, structs, variables. 

Use New as prefix for constructor functions. If a package has a single struct to be constructed, use New as the constructor function name.

Interfaces should be named with a verb + er prefix, e.g. `MessageBusConnecter`

###Nested code
Deeply nested code (> 3 levels) should be avoided. It breaks code readability and complicates understanding. 
Use early returns, functions and logical code separation to avoid nested code.

Let's consider typical cases of nested code and suggested solutions:

**Nested conditions**

        err := callZero()
        if err == nil {
            err = callOne()
            if err == nil {
                err = callThree()
                if err == nil {
                    doOtherWork()
                }
            }
        }
        
Solution using early returns:

        err := callZero()
        if err != nil {
            return err
        }
        
        err = callOne()
        if err != nil {
            return err
        }
        
        err = callThree()
        if err != nil {
            return err
        }
        doOtherWork()

**Nested loops**

        for _, itemsGroup := range items {
        	for _, itemsMap := range itemsGroup {
        		for _, item := range itemsMap {
        	        doSomeItemJob(item)
                }
            }
        }

Solution using functions:

        func processItems(items [][]map[string]Item) {
        	for _, itemsGroup := range items {
        		processItemsGroup(itemsGroup)
            }
        }
        
        func processItemsGroup(itemsGroup []map[string]Item){
        	for _, itemsMap := range itemsGroup {
        		processItemsMap(itemsMap)
            }
        }
        
        func processItemsMap(itemsMap map[string]Item){
        	for _, item := range itemsMap {
        		doSomeItemJob(item)
            }
        }

**Nested functions**

        func generateRandomText() string {
        	mutator := func(input string) string {
        		doubler := func (input string) string {
        			return input + input
                }
                cutter := func (input string) string {
        		    return input[0:len(input)-1]
                }
                return cutter(doubler(input))
            }
        	return mutator(“abc”)
        }
        
Refactored solution with logical code separation:


        type mutator func(input string) string

        func Generate(input string, mutators ...mutator) string {
        	for _, m := range mutators {
        	    input = m(input)
            }
            return input
        }
        
        doubler := func (input string) string {
            return input + input
        }
        
        cutter := func (input string) string {
            return input[0:len(input)-1]
        }
        
        ...
        return Generate(“abc”, doubler, cutter)


##References
- https://golang.org/doc/effective_go.html
- https://gist.github.com/wojteklu/73c6914cc446146b8b533c0988cf8d29 as summary of https://www.amazon.com/Clean-Code-Handbook-Software-Craftsmanship/dp/0132350882
- https://medium.com/mindorks/how-to-write-clean-code-lessons-learnt-from-the-clean-code-robert-c-martin-9ffc7aef870c
- https://golang.org/ref/spec









      



