# Building Sites with Sprout

Complete reference guide for building websites with Sprout static site generator.

## Table of Contents

1. [Core Concepts](#core-concepts)
2. [Directory Structure](#directory-structure)
3. [Content Files](#content-files)
4. [Templates](#templates)
5. [Styling and Assets](#styling-and-assets)
6. [URLs and Links](#urls-and-links)
7. [Building and Deployment](#building-and-deployment)
8. [Advanced Topics](#advanced-topics)
9. [Troubleshooting](#troubleshooting)

## Core Concepts

### What Sprout Does

Sprout converts Markdown files into HTML pages using HTML templates. It:
- Reads Markdown files from `content/`
- Applies HTML templates from `templates/`
- Generates static HTML files in `public/`
- Copies static assets (CSS, images) from `static/` to `public/`

### The Build Process

1. **Scan**: Finds all `.md` files in `content/`
2. **Parse**: Extracts front matter and converts Markdown to HTML
3. **Render**: Applies templates to create final HTML
4. **Output**: Writes HTML files to `public/` with pretty URLs
5. **Copy**: Copies static files from `static/` to `public/`

### Incremental Builds

Sprout tracks file changes and only rebuilds what's modified:
- Changed Markdown file → rebuilds only that page
- Changed template → rebuilds all pages (templates affect everything)
- Changed static file → copies only that file

This makes subsequent builds very fast.

## Directory Structure

### Required Directories

```
my-site/
├── content/          # Markdown source files (required)
├── templates/        # HTML templates (required)
├── static/          # CSS, images, JS (optional)
└── public/          # Generated site (auto-created)
```

### Directory Purposes

- **content/**: All your Markdown files go here. One `.md` file = one page.
- **templates/**: HTML templates that define page structure and styling.
- **static/**: Files copied as-is to `public/`. Use for CSS, images, JavaScript.
- **public/**: Generated HTML files. This is what you deploy. Don't edit directly.

### Configuration File

`sprout.toml` in the root directory configures paths and settings:

```toml
base_url = "http://localhost:1313"
pretty_urls = true
unsafe_html = false

[paths]
content = "content"
templates = "templates"
static = "static"
public = "public"
```

**Configuration Options:**

- `base_url` (string): Base URL for your site. Used in templates and for generating absolute URLs.
- `pretty_urls` (boolean): Generate clean URLs with trailing slashes (default: `true`).
- `unsafe_html` (boolean): Allow raw HTML in Markdown files (default: `false`). When `false`, HTML tags in Markdown are stripped for security. Set to `true` to enable raw HTML rendering.

## Content Files

### File Format

Content files are Markdown (`.md`) with optional front matter:

```markdown
+++
title = "Page Title"
slug = "page-slug"
layout = "page"
+++

# Page Title

Your content in Markdown format.
```

### Front Matter

Front matter is TOML between `+++` delimiters at the top of the file.

**Available Fields:**

- `title` (string): Page title. If omitted, uses first H1 heading.
- `slug` (string): URL slug. If omitted, uses filename without extension.
- `layout` (string): Template name. If omitted, uses `page`.

**Important Notes:**

- If you specify a `slug`, it overrides the directory structure in the URL.
- To preserve directory structure, omit the `slug` field.
- Front matter is optional - files without it work fine.

### Raw HTML in Markdown

By default, Sprout strips raw HTML from Markdown files for security. To enable raw HTML:

1. Set `unsafe_html = true` in your `sprout.toml`:

```toml
unsafe_html = true
```

2. Use HTML directly in your Markdown:

```markdown
<div class="custom-class">
    <p>This HTML will be rendered as-is.</p>
</div>
```

**Security Note**: Only enable `unsafe_html` if you trust the content of your Markdown files. Raw HTML can be used for XSS attacks if content comes from untrusted sources.

### File Naming

- Use lowercase filenames with hyphens: `my-page.md`
- `index.md` creates the root URL `/` or section root
- File extension must be `.md`

### Content Organization

**Flat Structure** (simple sites):

```
content/
├── index.md
├── about.md
└── contact.md
```

**Nested Structure** (organized sites):

```
content/
├── index.md
├── about.md
├── blog/
│   ├── index.md
│   ├── post-1.md
│   └── post-2.md
└── docs/
    ├── index.md
    ├── guide-1.md
    └── guide-2.md
```

## Templates

### Template System

Sprout uses Go's `html/template` package. Templates are HTML files with template syntax.

### Required Templates

- **base.html**: Base template that defines page structure (required)
- **page.html**: Default page template (recommended)

### Template Variables

Templates have access to these variables:

- `{{.Site.BaseURL}}`: Site base URL from config
- `{{.Page.Title}}`: Page title
- `{{.Page.Slug}}`: Page slug
- `{{.Page.RelPermalink}}`: Relative URL (e.g., `/about/`)
- `{{.Page.ContentHTML}}`: Rendered HTML content
- `{{.Page.SourcePath}}`: Source file path
- `{{.Page.Layout}}`: Layout name

### Base Template

`templates/base.html` defines the overall page structure:

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{.Page.Title}} - My Site</title>
    <link rel="stylesheet" href="/style.css">
</head>
<body>
    <header>
        <nav>
            <a href="/">Home</a>
            <a href="/about/">About</a>
        </nav>
    </header>
    <main>
        {{template "content" .}}
    </main>
    <footer>
        <p>&copy; 2024 My Site</p>
    </footer>
</body>
</html>
```

The `{{template "content" .}}` line includes content from other templates.

### Page Template

`templates/page.html` defines the default page content:

```html
{{define "content"}}
<article>
    <h1>{{.Page.Title}}</h1>
    <div class="content">
        {{.Page.ContentHTML}}
    </div>
</article>
{{end}}
```

The `{{define "content"}}` block is included by `base.html`.

### Custom Layouts

Create additional templates for different page types:

```html
<!-- templates/post.html -->
{{define "content"}}
<article class="post">
    <header>
        <h1>{{.Page.Title}}</h1>
    </header>
    <div class="post-content">
        {{.Page.ContentHTML}}
    </div>
</article>
{{end}}
```

Use it by setting `layout = "post"` in front matter.

### Template Functions

Sprout provides minimal template functions. Use Go template syntax:

- `{{if eq .Page.RelPermalink "/"}}active{{end}}`: Conditional rendering
- `{{.Page.Title | html}}`: HTML escaping (automatic for ContentHTML)

## Styling and Assets

### CSS Files

Place CSS files in `static/`:

```
static/
└── style.css
```

They're copied to `public/` during build:

```
public/
└── style.css
```

### Linking CSS

Reference CSS in your `base.html` template:

```html
<head>
    <link rel="stylesheet" href="/style.css">
</head>
```

Paths start with `/` and reference files in `public/` after build.

### CSS Organization

Organize CSS in subdirectories:

```
static/
├── css/
│   ├── main.css
│   └── components.css
└── style.css
```

Link them all:

```html
<link rel="stylesheet" href="/style.css">
<link rel="stylesheet" href="/css/main.css">
<link rel="stylesheet" href="/css/components.css">
```

### Images

Place images in `static/`:

```
static/
└── images/
    ├── logo.png
    └── photo.jpg
```

Reference in Markdown:

```markdown
![Logo](/images/logo.png)
```

Or in templates:

```html
<img src="/images/logo.png" alt="Logo">
```

### JavaScript

Place JavaScript files in `static/`:

```
static/
└── js/
    └── app.js
```

Include in templates:

```html
<script src="/js/app.js"></script>
```

### Asset Paths

All asset paths in templates and Markdown should:
- Start with `/` (root-relative)
- Reference the final location in `public/`
- Use forward slashes

Examples:
- `/style.css` (not `style.css` or `./style.css`)
- `/images/logo.png` (not `images/logo.png`)

## URLs and Links

### URL Generation

Sprout generates clean URLs with trailing slashes:

- `content/index.md` → `/`
- `content/about.md` → `/about/`
- `content/blog/post.md` → `/blog/post/`
- `content/docs/guide.md` → `/docs/guide/`

All URLs end with `/` and create `index.html` files in directories.

### Internal Links

Link to other pages using their permalinks:

```markdown
[About Page](/about/)
[Blog Post](/blog/my-post/)
[Home](/)
```

**Important**: Always use trailing slashes for internal links.

### Anchor Links

Link to sections within a page:

```markdown
## Section Title

Content here...

[Link to section](#section-title)
```

Markdown automatically generates IDs from headings:
- Lowercase
- Spaces become hyphens
- Special characters removed

### External Links

Standard markdown links:

```markdown
[External Site](https://example.com)
```

### Navigation in Templates

Build navigation using template conditionals:

```html
<nav>
    <a href="/" {{if eq .Page.RelPermalink "/"}}class="active"{{end}}>Home</a>
    <a href="/about/" {{if eq .Page.RelPermalink "/about/"}}class="active"{{end}}>About</a>
    <a href="/blog/" {{if eq .Page.RelPermalink "/blog/"}}class="active"{{end}}>Blog</a>
</nav>
```

## Building and Deployment

### Building

Generate the site:

```bash
sprout build
```

Output goes to `public/` directory.

### Clean Build

Remove all generated files and rebuild:

```bash
sprout build --clean
```

Useful when templates change or you want a fresh build.

### Development Server

Preview your site locally:

```bash
sprout serve
```

Visits `http://localhost:1313` in your browser.

### Live Reload

Enable automatic page reload during development:

```bash
sprout serve --livereload
```

Pages automatically reload when you save changes.

## Advanced Topics

### Unsafe HTML Configuration

The `unsafe_html` setting controls whether raw HTML in Markdown files is rendered or stripped:

- **`unsafe_html = false`** (default): All HTML tags in Markdown are removed. This is the safe default.
- **`unsafe_html = true`**: Raw HTML in Markdown is rendered as-is.

**When to enable:**

- You need custom HTML elements in your content
- You're embedding iframes, custom components, or complex layouts
- You trust all content sources

**When to keep disabled:**

- Content comes from user input or external sources
- You want maximum security
- You prefer pure Markdown without HTML

**Example usage:**

```markdown
+++
title = "Custom Layout"
+++

## My Content

<div class="custom-wrapper">
    <p>This paragraph is wrapped in a custom div.</p>
    <iframe src="https://example.com/embed"></iframe>
</div>
```

With `unsafe_html = true`, the HTML is rendered. With `unsafe_html = false`, only the text content appears.

### Subdirectories and Organization

Use subdirectories to organize content:

```
content/
├── index.md
├── about.md
└── blog/
    ├── index.md          # Blog section homepage
    ├── 2025/
    │   ├── post-1.md
    │   └── post-2.md
    └── 2024/
        └── old-post.md
```

This creates:
- `/` - Home
- `/about/` - About page
- `/blog/` - Blog index
- `/blog/2025/post-1/` - Blog post
- `/blog/2025/post-2/` - Blog post
- `/blog/2024/old-post/` - Blog post

### Index Pages in Subdirectories

Create `index.md` in subdirectories for section landing pages:

```markdown
+++
title = "Blog"
+++

# Blog

All my blog posts...
```

This creates `/blog/` as a dedicated page.

### Preserving Directory Structure

To preserve directory structure in URLs, **don't specify a slug**:

```markdown
+++
title = "My Post"
+++
```

This uses the file path: `content/blog/post.md` → `/blog/post/`

If you specify a slug, it overrides the path:

```markdown
+++
title = "My Post"
slug = "custom-url"
+++
```

This creates `/custom-url/` regardless of file location.

### Template Composition

Use `{{define}}` blocks for reusable components:

```html
<!-- templates/base.html -->
{{define "header"}}
<header>
    <h1>My Site</h1>
</header>
{{end}}

<body>
    {{template "header" .}}
    {{template "content" .}}
</body>
```

### Multiple Layouts

Create different layouts for different content types:

```html
<!-- templates/article.html -->
{{define "content"}}
<article class="article">
    <h1>{{.Page.Title}}</h1>
    {{.Page.ContentHTML}}
</article>
{{end}}
```

Use with `layout = "article"` in front matter.

## Troubleshooting

### Pages Not Appearing

**Problem**: Page doesn't show up after build.

**Solutions**:
- Check file has `.md` extension
- Verify file is in `content/` directory
- Run `sprout build` again
- Check for errors in front matter (must be valid TOML)
- Ensure templates directory exists

### Links Not Working

**Problem**: Internal links return 404.

**Solutions**:
- Use trailing slashes: `/about/` not `/about`
- Verify target page exists and was built
- Check permalink matches link URL
- Ensure target file has correct front matter

### Styles Not Loading

**Problem**: CSS doesn't apply.

**Solutions**:
- Verify CSS file is in `static/` directory
- Check link path in template starts with `/`
- Run `sprout build` after adding CSS
- Check browser console for 404 errors
- Verify file was copied to `public/`

### Template Errors

**Problem**: Template rendering fails.

**Solutions**:
- Check template syntax (valid HTML + Go template)
- Verify `base.html` exists
- Ensure `{{define "content"}}` blocks match
- Check for unclosed template tags
- Validate TOML in front matter

### Build Errors

**Problem**: Build fails with errors.

**Solutions**:
- Check error message for specific file
- Verify all required directories exist
- Ensure templates are valid HTML
- Check file permissions
- Try `sprout build --clean`

### Images Not Displaying

**Problem**: Images don't show up.

**Solutions**:
- Verify image is in `static/` directory
- Check path starts with `/` in markdown/templates
- Run `sprout build` after adding images
- Verify file was copied to `public/`
- Check file extension is correct

### Incremental Builds Not Working

**Problem**: Changes don't appear after rebuild.

**Solutions**:
- Try `sprout build --clean`
- Check `.sprout-cache.json` exists in `public/`
- Verify file timestamps are updating
- Delete cache file and rebuild

## Quick Reference

### File Locations

- Content: `content/*.md`
- Templates: `templates/*.html`
- Static assets: `static/**`
- Output: `public/**`

### URL Rules

- `content/index.md` → `/`
- `content/page.md` → `/page/`
- `content/section/page.md` → `/section/page/`
- Custom slug overrides path

### Template Variables

- `{{.Page.Title}}` - Page title
- `{{.Page.ContentHTML}}` - Rendered content
- `{{.Page.RelPermalink}}` - Page URL
- `{{.Site.BaseURL}}` - Site URL

### Commands

- `sprout init` - Initialize new site
- `sprout build` - Build site
- `sprout build --clean` - Clean build
- `sprout serve` - Development server
- `sprout serve --livereload` - Server with live reload
