---
backend: "NAME OF BACKEND"
headline: "SHORT DESCRIPTION"
package: ""
features: []
project_url: ""
---
# {{backend}}

[![API docs](https://godoc.org/github.com/reddec/storages/{{package}}?status.svg)](http://godoc.org/github.com/reddec/storages/{{package}})

* **import:** `github.com/reddec/storages/{{package}}`
* [{{backend}} project]({{project_url}})

DESCRIPTION


## Usage

**Example**

```go
// TODO: code
```


## Features

{% for feature in features %}
{% include feature_{{feature}}.md %}
{% endfor %}
