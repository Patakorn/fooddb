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

// Category images
{{ template "imagePath" .Categories }}

// Ingredient images
{{ template "imagePath" .Ingredients }}

export const Languages = ["en", "fr"] as const
export type Language = typeof Languages[number]

export const Units = ["g", "Kg", "mL", "cL", "L", "each"] as const
export type Unit = typeof Units[number]

export type IngredientCategory = {
  id: string
  name: Partial<Record<Language, string>>
  img?: string
}

export type Ingredient = {
  id: string
  name: Record<Language, string>
  img?: string
  categories: string[]
  units: Unit[]
}

const categoryCatalog: Record<string, IngredientCategory> = {
{{- range $catName, $cat := .Categories }}
  {{ $catName }}: {
    id: "{{ $catName }}",
    name: {
    {{- range $k, $v := $cat.Name }}
      {{ $k }}: "{{ $v }}",
    {{- end }}
    },
    {{- if $cat.ImagePath }}
    img: {{ $catName }}_img,
    {{- end }}
  },
{{- end }}
}

const ingredientCatalog: Record<string, Ingredient> = {
{{- range $ingName, $ing := .Ingredients }}
  {{ $ingName }}: {
    id: "{{ $ingName }}",
    name: {
    {{- range $k, $v := $ing.Name }}
      {{ $k }}: "{{ $v }}",
    {{- end }}
    },
    {{- if $ing.ImagePath }}
    img: {{ $ingName }}_img,
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

export const FoodDB = {
  categories: categoryCatalog,
  ingredients: ingredientCatalog,
}
export default FoodDB
`
