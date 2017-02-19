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
` + `/* Global styling. */

body {
    font-family: "Helvetica Neue", Helvetica, "Segoe UI", Arial, freesans, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Microsoft Yahei", FZHei-B01, "WenQuanYi Micro Hei", SimHei;
    font-size: 16px;
}

#blog-title, h1, h2, h3 {
    font-family: Georgia, Times, Serif, FZShuSong-Z01, SimSun;
}

/* Layout primitives. */

.card {
    /*background-color: white;*/
    margin: 0;
}

#navbar-container {
    top: 0;
    width: 100%;
    color: white;
    background-color: #1a1a1a;
}

@media screen and (min-width: 1024px) {
    .card {
        width: 1024px;
        margin: 0 auto 0 auto;
        padding: 32px 0 20px;
    }
    #navbar-container {
        border-bottom: 0;
    }
    /* Like .card, with no margin on top or bottom. */
    #navbar {
        width: 1000px;
        margin: 0 auto;
    }
}

.nav-link {
    color: black;
    display: inline;
    text-decoration: none;
    line-height: 1.4em;
}

.clear {
    clear: both;
}

img {
    max-width: 100%;
}

/* Global header card. */

#blog-title {
    padding: 2% 4% 1%;
    font-size: 2em;
    font-weight: bold;
    /* color: #f0f0f0; */
    color: white;
}

.card-splitter {
    height: 1px;
    width: 100%;
    background-color: #667;
}

#category-list {
    display: block;
    padding: 1.6% 4%;
}

#category-list > li {
    list-style: none;
    display: inline-block;
    margin: 4px 0;
}

#category-list > li > .nav-link {
    color: white;
    padding: 4px 10px 4px;
}

#category-list > li > .nav-link:hover {
    background-color: #333;
}

#category-list > li.current > .nav-link {
    color: black;
    background-color: white;
}

/* Article content. */

.article-header {
    padding: 3.2% 4% 0;
}

.article-content {
    padding: 4%;
}

.article p, .article pre {
    line-height: 1.6em;
    margin-top: 0.5em;
    margin-bottom: 1em;
}

.article code, .article pre {
    background-color: #f0f0f0;
    border-radius: 3px;
    font-family: Consolas, "Liberation Mono", Menlo, Courier, monospace;
    font-size: 85%;
    padding: 0.2em;
}

.article pre {
    overflow: auto;
}

/* Only use headers up to h3 */
.article h1, .article h2, .article h3 {
    /*font-family: Arial, Sans Serif, Microsoft Yahei, FZHei-B01, WenQuanYi Micro Hei, SimHei;*/
    line-height: 1.6em;
    margin-top: 0.6em;
    margin-bottom: 1em;
    font-weight: bold;
}

.article h2 {
    border-bottom: 1px solid #ddd;
}

.article h1 {
    font-size: 1.4em;
}

.article h2 {
    font-size: 1.2em;
}

.article h3 {
    font-size: 1em;
}

.article ul, .article ol {
    margin-left: 1em;
}

/* Category content. */

.category-prelude {
    padding: 4% 5.8% 0;
    margin-bottom: -20px;
}

.article-list {
    padding: 4%;
}

.article-list > li {
    list-style: square inside;
    padding: 4px;
}

.article-list > li:hover {
    background-color: #c0c0c0;
}

.article-list > li > .nav-link {
    border-bottom: 1px solid black;
}

.article-timestamp {
    float: right;
    display: inline-block;
}

/* Utilities usable in articles. */

hr {
    clear: both;
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

/* vi: se ts=4 sts=4 sw=4: */
`

const baseTemplText = `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  {{ if is "homepage" }}
    <link rel="alternate" type="application/atom+xml"
          href="{{ rootURL }}/feed.atom">
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
  <div id="navbar-container"> <div id="navbar">
    <div id="blog-title">
      {{ .BlogTitle }}
    </div>
    <ul id="category-list">
      {{ $homepageTitle := .HomepageTitle }}
      {{ $curcat := .Category }}
      <li class="{{ if eq $curcat "homepage" }}current{{ end }}">
        <a href="{{ rootURL }}" class="nav-link">
          {{ $homepageTitle }}
        </a>
      </li>
      {{ range $info := .Categories }}
        <li class="{{ if eq $curcat $info.Name}}current{{ end }}">
          <a href="{{ rootURL }}/{{ $info.Name }}" class="nav-link">
            {{ $info.Title }}
          </a>
        </li>
      {{ end }}
    </ul>
  </div> </div>

  {{/*
    The reference to "content" is a free one and has to be fixed elsewhere.
    The *-content templates defined below are intended to be used for this.

    For instance, by adding the following code, this whole template file will
    function as the template for articles:

        {{ define "content" }} {{ template "article-content" . }} {{ end }}

    This snippet can be generated by contentIs("article").
  */}}
  {{ template "content" . }}
</body>
</html>

{{ define "article-content" }}
  <div class="card">
    <article class="article">
      {{ if not .IsHomepage }}
        <div class="article-header">
          <div class="timestamp"> {{ .Timestamp }} </div>
          <h1> {{ .Title }} </h1>
          <div class="clear"></div>
        </div>
        <div class="card-splitter"></div>
      {{ end }}
      <div class="article-content">
        {{ .Content }}
      </div>
      <div class="clear"></div>
    </article>
  </div>
{{ end }}

{{ define "category-content" }}
  {{ $category := .Category }}
  <div class="card">
    {{ if ne .Prelude "" }}
      <div class="category-prelude article">
        <article class="article">
          {{ .Prelude }}
        </article>
      </div>
    {{ end }}
    <ul class="article-list">
      {{ range $article := .Articles }}
        <li>
          <a href="{{ rootURL }}/{{ $category }}/{{ $article.Name }}.html"
             class="nav-link">{{ $article.Title }}</a>
          <span class="article-timestamp">
            {{ $article.Timestamp }}
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
