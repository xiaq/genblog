package main

const defaultCSS = `/*
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

html {
  /* Prevent font scaling in landscape while allowing user zoom */
  -webkit-text-size-adjust: 100%;
}

body {
    font-family: "Helvetica Neue", Helvetica, "Segoe UI", Arial, freesans, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Microsoft Yahei", FZHei-B01, "WenQuanYi Micro Hei", SimHei;
    font-size: 16px;
    background-color: white;
}

/*
#blog-title, h1, h2, h3 {
    font-family: Georgia, Times, Serif, FZShuSong-Z01, SimSun;
}
*/

/* Layout primitives. */

#navbar-container {
    top: 0;
    width: 1280px;
    max-width: 100%;
    margin: auto;
    color: white;
    background-color: #1a1a1a;
}

.card {
    background-color: white;
}

.card, #navbar {
    margin: 0 auto;
    max-width: 1024px;
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

a {
    color: #6666ff;
}

/* Global header card. */

#navbar {
    padding-bottom: 16px;
}

#blog-title {
    font-size: 2em;
    color: white;
    display: inline-block;
    margin: 16px 2% 0 4%;
}

.card-splitter {
    height: 1px;
    width: 100%;
    background-color: #667;
}

#category-list {
    display: inline-block;
    /* The top margin is used when the category list is not on the same line
     * as the title. Otherwise it is not important. */
    margin-top: 16px;
    /* Used (together with margin-right of #blog-title) to separate category
     * list and blog title when they are on the same line. Otherwise, used as
     * spacing from the page edge. This needs to be the same as that of
     * #blog-title to keep them left-aligned. */
    margin-left: 2%;
}

#category-list > li {
    list-style: none;
    display: inline-block;
    /* This margin is important for the correct spacing when the category list
     * runs multiple lines. */
    margin-top: 10px;
}

#category-list > li > .nav-link {
    color: white;
    position: relative;
    top: -4px; /* Cancel padding-bottom to align with baseline. */
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

.timestamp {
    margin-bottom: 0.6em;
}

.article-header {
    padding: 2.4% 4% 1.6%;
}

.article-content {
    padding: 4%;
}

.toc {
    margin-left: -0.8em;
}

.toc li {
    list-style: none;
}

.toc a {
    color: black;
}

.article p, .article ul, .article pre {
    margin-bottom: 16px;
}

.article p {
    line-height: 1.5;
}

.article li {
    margin-top: 0.5em;
    margin-bottom: 0.5em;
}

.article li > p {
    margin-top: 1em;
    margin-bottom: 1em;
}

.article pre {
    padding: 1em;
    line-height: 1.45;
    overflow: auto;
}

.article p code {
    padding: 0.2em 0;
}

.article p code::before, .article p code::after {
    letter-spacing: -0.2em;
    content: "\00a0";
}

code, pre {
    font-family: Consolas, Menlo, "Bitstream Vera Sans Mono", "DejaVu Sans Mono", monospace;
}

.article code, .article pre {
    background-color: #f0f0f0;
    border-radius: 3px;
}

.article code {
    font-size: 85%;
}

/* Only use headers up to h2 */
.article-header h1, .article-content h1, .article-content h2 {
    line-height: 1.25;
}

.article-header h1 {
    font-size: 2em;
}

.article-content h1, .article-content h2 {
    margin-top: 24px;
    margin-bottom: 16px;
}

.article-content h1 {
    font-size: 2em;
    padding-bottom: 0.3em;
    border-bottom: 1px solid #eee;
    font-weight: 600;
}

.article-content h2 {
    font-size: 1.5em;
    padding-bottom: 0.3em;
    border-bottom: 1px solid #eee;
    font-weight: 600;
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
    text-shadow: 0px 0px 2px #667;
    display: inline-block;
    position: relative;
    top: -0.5em;
    padding: 0 0.25em;
    font-size: 1.1em;
    color: black;
    background-color: white;
}

/* vi: se ts=4 sts=4 sw=4: */
`

const defaultTemplate = `<!doctype html>
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
    <div class="clear"></div>
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
