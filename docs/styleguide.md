# Style guide

Go style conventions for this module. Adapted from the shared trancecode Go style guide.

## Errors

### General guidelines

* Never log and return; one or the other
* Never just return `err`, always provide additional context. If there is none to add, comment what is already there (e.g. `// os.PathError already includes operation and filename`)

### Formatting

Error messages should be phrased as `<context>: <reason>` where `<reason>` is generally the underlying error message.

The context should indicate the action being attempted.

Example:
```Go
if err := doSomething(); err != nil {
	return fmt.Errorf("doing something: %w", err)
}
```

### Providing context

* Context includes things like loop iterations and computed values the caller doesn't know or the reader might need
* Context includes what the code block is trying to do, not internals like function names
* Context must uniquely identify the code path when there could be multiple error returns
* Don't hesitate to use %T when dealing with unknown types
* Always use %q for strings we can't guarantee are clean, non-empty strings
* Include all context that the caller doesn't have, omit most context the caller does have
* Don't start with `failed to` or `error` except when you are logging

### Examples

* `looking up modifier for ability 'STR': no such ability 'STR'`

## Panic

Prefer `panic()` to `log.Fatalf()`. Use panic() for unrecoverable errors that should terminate the program immediately, and for invalid arguments that indicate a programming error (e.g. an unspecified enum value).

## Naming

* Use Java camel case convention for acronyms. e.g. use `XpThreshold` instead of `XPThreshold`, and `Id` instead of `ID`.
* Use the idiomatic Go naming convention for everything else.

## Imports

* **Group imports:**
    * Standard library imports
    * Related third party imports
    * Local application/library specific imports
* **Import order within groups:** Sort alphabetically.

## Enumerations

When defining enumerations using `iota`, always start with a "None", "Unspecified", or "Invalid" value as the zero value. This makes it explicit when enum fields are uninitialized or invalid.

**Preferred pattern:**
```go
type Size int
const (
    SizeUnspecified Size = iota  // default/uninitialized value
    SizeTiny
    SizeSmall
    SizeMedium
    SizeLarge
)
```

**Benefits:**
* **Explicit validation**: Can detect uninitialized enum fields (`if size == SizeUnspecified`)
* **Debugging**: Clear indication when enum values haven't been set properly
* **Defensive programming**: Prevents subtle bugs from relying on implicit zero values
* **Code clarity**: Makes intent explicit rather than relying on implicit behavior

This pattern should be applied to all new enumerations in the codebase.

### Enumerations: validation and error handling

When working with enumerations, always validate that enum values are not in their "None"/"Unspecified" state before using them:

**Preferred:**
```go
func sizeFactor(size Size) int {
    if size == SizeUnspecified {
        panic("size must be specified")
    }
    // proceed
}
```

**Also good for non-critical paths:**
```go
func process(size Size) {
    if size == SizeUnspecified {
        return // skip invalid values
    }
    // process
}
```

This ensures that enum fields are always explicitly set and prevents subtle bugs from uninitialized values.

## Comments

* **Write clear and concise comments:** Explain the "why" behind the code, not just the "what".
* **Comment sparingly:** Well-written code should be self-documenting where possible.
* **Use complete sentences:** Start comments with a capital letter and use proper punctuation.
* **Consider context:** Place comments where they are most helpful for understanding the code, often above or to the right of the relevant code.
* **Be concise:** Keep comments short and to the point, avoiding unnecessary verbosity.
* **Be consistent:** Follow the established conventions within the project and the broader Go community.
* **Document gotchas:** Explain any potential pitfalls or unusual behavior of the code.

## Documentation for packages, types and functions

* Document all exported types and functions.
* Start with the name: doc comments should generally begin with the name of the thing they're documenting (function, type, variable, etc.).
* Explain purpose, not implementation: focus on what the code is supposed to do, not the specific code logic.
* Keep the documentation concise and straightforward. The goal is to convey as much information as possible without using up too much of the context window. Explain the overall purpose, give details about corner cases.
* When describing a function: no need to describe each argument on its own. Just describe the function and how it uses the arguments.
* Use proper formatting: utilize bullet points, code blocks (using tabs), and links where appropriate, leveraging Go's Markdown-like formatting.

### Struct field documentation

* **Document all exported struct fields:** Each exported field should have a comment explaining its purpose and usage.
* **Place field comments above the field:** Use the line above the field declaration for documentation, not inline comments.
* **Start with the field name:** Begin the comment with the field name followed by a description.
* **Be concise and specific:** Focus on the field's role and include units, ranges, or constraints where applicable.

**Preferred pattern:**
```go
// AbilityScores holds the six core ability scores of a creature.
type AbilityScores struct {
    // Strength measures physical power, affecting melee attacks and carrying capacity.
    Strength int

    // Dexterity measures agility, affecting armor class and initiative.
    Dexterity int
}
```

**Avoid:**
```go
type AbilityScores struct {
    Strength  int // strength
    Dexterity int // dexterity
}
```

#### Specific examples

Example 1:
* Instead of: `// Increment the counter variable.`
* Prefer: `// Counter increments the value of a counter.`

Example 2:
* Instead of: `// This function calculates the sum of two numbers.`
* Prefer: `// Sum calculates the sum of two integers.`

## Dos and don'ts

### Simple conditionals: avoid `else` blocks where applicable

Instead of:
```go
if condition {
  // do something
} else {
  return
}
```

Prefer:
```go
if !condition {
  return
}

// do something
```
