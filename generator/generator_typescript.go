package generator

const TypeScriptImageDir string = "images"

const TypeScriptModuleTemplate string = `
// Generated from FoodDB
// https://github.com/fuhrmannb/fooddb

declare module '*.jpg' {
  const src: string;
  export default src;
}

declare module '*.jpeg' {
  const src: string;
  export default src;
}

declare module '*.png' {
  const src: string;
  export default src;
}
`

const TypeScriptTemplate string = `
// Generated from FoodDB
// https://github.com/fuhrmannb/fooddb

{{- define "imagePath" -}}
  {{- range $name, $cat := . }}
    {{- if $cat.ImagePath -}}
      import {{ $name }}_img from "./images/{{ $name }}.jpg"
    {{- end }}
  {{- end }}
{{- end }}
{{ template "imagePath" .Categories }}
{{ template "imagePath" .Ingredients }}

type Languages = "en" | "fr"

export type IngredientCategory = {
  name: Partial<Record<Languages, string>>
  image?: string
}

type Units = "g" | "Kg" | "L" | "each"

export type Ingredient = {
  name: Record<Languages, string>
  image?: string
  categories: string[]
  units: Units[]
}

export const IngredientCategoryCatalog: Record<string, IngredientCategory> = {
{{- range $catName, $cat := .Categories }}
  {{ $catName }}: {
    name: {
    {{- range $k, $v := $cat.Name }}
      {{ $k }}: "{{ $v }}",
    {{- end }}
    },
    {{- if $cat.ImagePath }}
      image: {{ $catName }}_img,
    {{- end }}
  },
{{- end }}
}

export const IngredientCatalog: Record<string, Ingredient> = {
{{- range $ingName, $ing := .Ingredients }}
  {{ $ingName }}: {
    name: {
    {{- range $k, $v := $ing.Name }}
      {{ $k }}: "{{ $v }}",
    {{- end }}
    },
    {{- if $ing.ImagePath }}
      image: {{ $ingName }}_img,
    {{- end }}
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
