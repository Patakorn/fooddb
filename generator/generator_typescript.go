package generator

const TypeScriptTemplate string = `
type Languages = "en" | "fr"

export type IngredientCategorie = {
  name: Partial<Record<Languages, string>>,
}

type Units = "g" | "Kg" | "L" | "each"

export type Ingredient = {
  name: Record<Languages, string>,
  image: string
  categories: string[]
  units: Units[]
}

const IngredientCategoriesCatalog: Record<string, IngredientCategorie> = {
{{- range $catName, $cat := .Categories }}
  {{ $catName }}: {
    name: {
    {{- range $k, $v := $cat.Name }}
      {{ $k }}: "{{ $v }}",
    {{- end }}
    },
  },
{{- end }}
}

const IngredientCatalog: Record<string, Ingredient> = {
{{- range $ingName, $ing := .Ingredients }}
  {{ $ingName }}: {
    name: {
    {{- range $k, $v := $ing.Name }}
      {{ $k }}: "{{ $v }}",
    {{- end }}
    },
    image: "{{ $ing.Image }}",
    categories: [
    {{- range $ing.Categories }}
      "{{ . }}",
    {{- end }}
    ],
    units: [
    {{- range $ing.Units }}
        "{{ . }}",
    {{- end }}
    ],
  },
  {{- end }}
}
`
