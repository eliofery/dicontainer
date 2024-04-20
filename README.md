# dicontainer : 🧶 is a simple dependency injection.

DI container provides a seamless way to wire up your application's components, making
it easy to manage dependencies and promote better code organization.

## Installation

```bash
go get github.com/eliofery/dicontainer
```

## Usage

```go
type FooStruct struct {}
type BarStruct struct {}
type BazStruct struct {val int}

di := dicontainer.New()

err := di.Set(
    func(bar BarStruct) *FooStruct {
        // Do something.
        _ = bar
        
        return &FooStruct{}
    },
    
    func() BarStruct {
        return BarStruct{}
    },
    
    func(foo *FooStruct) (BazStruct, error) {
        if rand.IntN(2) == 1 {
            return BazStruct{}, errors.New("baz error")
        }
        
        // Do something.
        _ = foo
        
        return BazStruct{123}, nil
    },
)

if err != nil {
    panic(err)
}

baz := di.Get("BazStruct").(BazStruct)

fmt.Println(baz)
```

## Inspiration

I was inspired to write this package by the [uber-go/fx](https://github.com/uber-go/fx).
