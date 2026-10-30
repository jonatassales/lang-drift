# Prints, once, each tracked path whose last segment is not kebab-case, for
# drift quality coverage. Conventional names, Next.js route segments and the forms
# the Go toolchain requires pass (.claude/rules/01-filesystem.rule.md).
{
	path = ""
	for (part = 1; part <= NF; part++) {
		name = $part
		path = path (part > 1 ? "/" : "") name
		if (name ~ /^(README|AGENTS|CLAUDE|SKILL)(\.[a-z-]+)?\.md(\.tmpl)?$/) continue
		if (name ~ /^(LICENSE|Dockerfile(\.[a-z]+)?)$/) continue
		sub(/^\./, "", name)
		gsub(/[\[\]()@]|\.\.\./, "", name)
		sub(/_test\.go$/, ".go", name)
		sub(/_(linux|darwin|windows|amd64|arm64)\.go$/, ".go", name)
		if (name ~ /^([a-z0-9]+(-[a-z0-9]+)*)?(\.[a-z0-9]+(-[a-z0-9]+)*)*$/) continue
		if (path in seen) continue
		seen[path] = 1
		print path
	}
}
