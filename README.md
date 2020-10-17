# Validator

Validator is a validate tool, it can help you

- check message parameter's type
- senitize message, from string to other type.

It's useful when you have to receive third party messages, but you're not sure its type

## Install

```
go get github.com/ken00535/validator
```

## Usage

For using this package, your message need to implement `Payload` interface

For example

```go
type message struct {
	msg   map[string]string
	cache map[string]interface{}
}

func (m *message) GetCache() map[string]interface{} {
	if m.cache == nil {
		m.cache = make(map[string]interface{})
	}
	return m.cache
}

func (m *message) SetCache(input map[string]interface{}) {
	m.cache = input
}

func (m *message) GetParam(field string) (val string, exist bool) {
	v, ok := m.msg[field]
	return v, ok
}
```

### Check

```go
payload := &message{msg: map[string]string{}}
payload.msg["age"] = "18"
Check(payload).Params("score").IsFloat()
errs, absence := ValidateResult(tc.dataReq)
```

### Sanitize

For using sanitize, your struct that carry message should add `vld` to field tag

```go
type person struct {
	Name    string  `vld:"name"`
	Gender  string  `vld:"gender"`
	Age     int     `vld:"age"`
	Score   int     `vld:"score"`
	Weight  float64 `vld:"w"`
	IsAlive bool    `vld:"alive"`
}
```

And you can call validator at project

```go
payload := &message{msg: map[string]string{}}
payload.msg["age"] = "18"
player := &person{Name: "ken"}
validator.Assign(payload).Struct(&player)
validator.Sanitize(payload).Params("age").ToInt()

fmt.Println(player.age)
// -> 18
```

## Error Handling

```go
payload := &message{msg: map[string]string{}}
payload.msg["age"] = "18"
Check(payload).Params("score").IsFloat()
errs, absence := ValidateResult(tc.dataReq)
if len(errs) > 0 {
	if err, ok := errs[0].(validator.Error); ok {
		if err.IsNotExist() {
			// do something
		}
	}
	return err
}
```
