package main

const css = `/*
Copied from:  http://code.google.com/p/html5resetcss/

html5doctor.com Reset Stylesheet
v1.6.1
Last Updated: 2010-09-17
Author: Richard Clark - http: //richclarkdesign.com
Twitter: @rich_clark
*/
html, body, div, span, object, iframe, h1, h2, h3, h4, h5, h6, p,
blockquote, pre, abbr, address, cite, code, del, dfn, em, img, ins,
kbd, q, samp, small, strong, sub, sup, var, b, i, dl, dt, dd, ol, ul, li,
fieldset, form, label, legend, table, caption, tbody, tfoot, thead, tr, th, td,
article, aside, canvas, details, figcaption, figure, footer, header, hgroup,
menu, nav, section, summary, time, mark, audio,
video { margin: 0; padding: 0; border: 0; outline: 0;
            font-size: 100%; vertical-align: baseline; background: transparent; }
body { line-height: 1; }
article, aside, details, figcaption, figure, footer, header, hgroup, menu, nav,
section { display: block; }
nav ul { list-style: none; }
blockquote, q { quotes: none; }
blockquote:before, blockquote:after, q:before, q:after { content: ''; content: none; }
a { margin: 0; padding: 0; font-size: 100%; vertical-align: baseline; background: transparent; }
ins { background-color: #ff9; color: #000; text-decoration: none; }
mark { background-color: #ff9; color: #000; font-style: italic; font-weight: bold; }
del { text-decoration: line-through; }
abbr[title], dfn[title] { border-bottom: 1px dotted; cursor: help; }
table { border-collapse: collapse; border-spacing: 0; }
hr { display: block; height: 1px; border: 0; border-top: 1px solid #cccccc; margin: 1em 0; padding: 0; }
input, select { vertical-align: middle; }
/* end HTML5 reset */
` + `body {
    background-color: #e5e5e5;
    font-size: 20px;
    font-family: Georgia, Times, Serif, FZShuSong-Z01, SimSun;
}

.card {
    background-color: white;
    padding: 24px;
    margin-top: 16px;
}

.card.noted {
    margin-top: 0;
}

.card-notes {
    margin-top: 16px;
}

@media screen and (min-width: 768px) {
    .card {
        margin: 32px auto 20px auto;
        padding: 48px;
        width: 672px;
        box-shadow: 2px 1px 3px #cce;
    }

    .card-notes {
        margin: 32px auto 0 auto;
        width: 672px;
    }
}

.card#header {
    font-size: 2em;
}

.category-list > li {
    list-style: none;
    display: inline-block;
    padding: 8px 10px 4px 10px;
}

.category-list > li.current {
    background-color: white;
}

.nav-link {
    color: black;
    display: inline-block;
    text-decoration: none;
    padding-bottom: 1px;
    margin-bottom: 2px;
    border-bottom: 2px solid #bbb;
}

.nav-link.current, .nav-link.current:hover {
    border-bottom-color: black;
}

.nav-link:hover {
    border-bottom-color: black;
}

.article p {
	line-height: 1.2;
	margin-bottom: 0.7em;
}

/* Only use headers up to h3 */
.article h1, .article h2, .article h3 {
    font-family: Arial, Sans Serif, FZHei-B01, WenQuanYi Micro Hei, SimHei;
    line-height: 1.5;
    margin-top: 0.6em;
    margin-bottom: 0.6em;
    font-weight: normal;
    font-variant: small-caps;
}

.article h1 {
    font-size: 1.4em;
}

.article h2 {
    font-size: 1.2em;
}

.article h2:before {
    content: "# ";
}

.article h3 {
    font-size: 1em;
}

.article h3:before {
    content: "## ";
}

.article ul, .article ol {
    margin-left: 1em;
}

.article-list > li {
    list-style: square inside;
    padding-bottom: 2px;
    margin-bottom: 0.7em;
}

.article-meta {
    float: right;
    display: inline-block;
}

.article-meta.header {
    margin-bottom: 0.7em;
}

.clear, hr {
    clear: both;
}

hr {
    border-color: #aaa;
    text-align: center;
}

hr:after {
    content: "❧";
    color: black;
    text-shadow: 0px 0px 2px #667;
    display: inline-block;
    position: relative;
    top: -0.5em;
    padding: 0 0.25em;
    font-size: 1.1em;
    background-color: white;
}
`

const baseTemplText = `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width; initial-scale=1.0; maximum-scale=1.0; user-scalable=0;" />
  {{ if is "homepage" }}
    <link rel="alternate" type="application/atom+xml"
          href="{{ .RootURL }}/feed.atom">
  {{ end }}

  <title>
    {{ if is "homepage" }}
      {{ .BlogTitle }}
    {{ else if is "category" }}
      {{ index .CategoryMap .Category }}
    {{ else }}
      {{ .Title }}
    {{ end }}
  </title>
  <style> {{ .CSS }} </style>
</head>
<body>
  <div class="card" id="header">
    <a href="{{ homepageURL }}" class="nav-link">{{ .BlogTitle }}</a>
  </div>
  {{ template "content" . }}
</body>
</html>

{{ define "category-list" }}
  {{ $cat := .Category }}
  <ul class="category-list">
    {{ range $info := .Categories }}
      <li class="{{ if eq $cat $info.Name }}current{{ end }}">
        <a href="{{ categoryURL $info.Name }}" class="nav-link {{ if and (eq $cat $info.Name) (is "category") }} current {{ end }}">
          {{ $info.Title }}
        </a>
      </li>
    {{ end }}
  </ul>
{{ end }}

{{ define "homepage-content" }}
  {{ range .Articles }}
    {{ template "article-content" . }}
  {{ end }}
{{ end }}

{{ define "article-content" }}
  <div class="card-notes">
    {{ template "category-list" . }}
  </div>
  <div class="card noted">
    <article class="article">
      <h1>
        <a href="{{ articleURL .Category .Name }}"
           class="nav-link {{ if is "article" }} current {{ end }}">
          {{ .Title }}
        </a>
      </h1>
      <div class="clear"></div>
      <span class="article-meta header">
        {{ index .CategoryMap .Category }} · {{ .Timestamp }}</span>
      <div class="clear"></div>
      {{ .Content }}
      <div class="clear"></div>
    </article>
  </div>
{{ end }}

{{ define "category-content" }}
  <div class="card-notes">
    {{ template "category-list" . }}
  </div>
  <div class="card noted">
    <ul class="article-list">
      {{ $catMap := .CategoryMap }}
      {{ range $info := .Articles }}
        <li>
          <a href="{{ articleURL $info.Category $info.Name }}"
             class="nav-link">{{ $info.Title }}</a>
          <span class="article-meta">
          {{ if not (is "category") }}
            {{ index $catMap $info.Category }} ·
          {{ end }}
          {{ $info.Timestamp }}
          </span>
          <div class="clear"></div>
        </li>
      {{ end }}
    </ul>
  </div>
{{ end }}

<!-- vi: se sw=2 ts=2 sts=2 et: -->
`

const feedTemplText = `<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
	<title>{{ .BlogTitle }}</title>
	<link href="{{ .RootURL }}"/>
	<link rel="self" href="{{ .RootURL }}/feed.atom"/>
	<updated>{{ .LastModified }}</updated>
	<id>{{ .RootURL }}/</id>

	{{ $rootURL := .RootURL }}
	{{ $author := .Author }}
	{{ range $info := .Articles}}
	<entry>
		<title>{{ $info.Title }}</title>
		{{ $link := print $rootURL "/" $info.Category "/" $info.Name ".html" }}
		<link rel="alternate" href="{{ $link }}"/>
		<id>{{ $link }}</id>
		<updated>{{ $info.LastModified }}</updated>
		<author><name>{{ $author }}</name></author>
		<content type="html">{{ $info.Content | html }}</content>
	</entry>
	{{ end }}
</feed>
`
