# Prints each tracked file under an app's src/app/atoms that is not an atom
# module, its test or its types (.claude/rules/07-jotai.rule.md, rule 3: one
# concern per <concern>.atom.ts, re-exported from the @state barrel).
$0 ~ /^apps\/[^\/]+\/src\/app\/atoms\// && $0 !~ /\.(atom|test|unit|types)\.ts$/ {
  print
}
