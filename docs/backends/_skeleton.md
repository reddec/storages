---
backend: "NAME OF BACKEND"
headline: "SHORT DESCRIPTION"
package: ""
features: []
project_url: ""
---
# {{page.backend}}

[![API docs](https://godoc.org/github.com/reddec/storages/{{page.package}}?status.svg)](http://godoc.org/github.com/reddec/storages/{{page.package}})

* **import:** `github.com/reddec/storages/{{page.package}}`
* [{{page.backend}} project]({{page.project_url}})

DESCRIPTION


## Usage

**Example**

```go
// TODO: code
```


## Features

{% for feature in page.features%}
{% include feature_{{feature}}.md %}
{% endfor %}
