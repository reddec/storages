{% assign featuresNum = include.page.features | size %}
{% if featuresNum > 0 %}
## Features

{% for feature in include.page.features%}
{% include feature_{{feature}}.md %}
{% endfor %}
{%endif%}
