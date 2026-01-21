package generator

// modelTemplate 标准Java Bean Model模板
const modelTemplate = `package {{.Package}};

import java.io.Serializable;
{{range .Imports}}import {{.}};
{{end}}
{{if .UseJsonProperty}}import com.fasterxml.jackson.annotation.JsonProperty;
{{end}}
{{if .TableComment}}/**
 * {{.TableComment}}
 */
{{end}}public class {{.ClassName}} implements Serializable {
    private static final long serialVersionUID = 1L;
{{range .Fields}}
{{if .Comment}}    /** {{.Comment}} */
{{end}}{{if $.UseJsonProperty}}{{if $.JsonPropertyUpperCase}}    @JsonProperty("{{title .ColumnName}}")
{{else}}    @JsonProperty("{{.ColumnName}}")
{{end}}{{end}}    private {{.FieldType}} {{.FieldName}};
{{end}}
{{if .NeedConstructors}}    public {{.ClassName}}() {}

    public {{.ClassName}}({{range $i, $e := .Fields}}{{if $i}}, {{end}}{{$e.FieldType}} {{$e.FieldName}}{{end}}) {
        {{range .Fields}}this.{{.FieldName}} = {{.FieldName}};
        {{end}}
    }
{{end}}
{{range .Fields}}
    public {{.FieldType}} get{{title .FieldName}}() {
        return {{.FieldName}};
    }

    public void set{{title .FieldName}}({{.FieldType}} {{.FieldName}}) {
        this.{{.FieldName}} = {{.FieldName}};
    }
{{end}}
{{if .NeedToStringHashcodeEquals}}    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        {{.ClassName}} that = ({{.ClassName}}) o;
        return {{range $i, $e := .Fields}}{{if $i}} && {{end}}Objects.equals({{.FieldName}}, that.{{.FieldName}}){{end}};
    }

    @Override
    public int hashCode() {
        return Objects.hash({{range $i, $e := .Fields}}{{if $i}}, {{end}}{{.FieldName}}{{end}});
    }

    @Override
    public String toString() {
        StringBuilder sb = new StringBuilder();
        sb.append("{{.ClassName}}{");
        {{range $i, $e := .Fields}}{{if $i}}sb.append(", ");{{end}}sb.append("{{.FieldName}}=").append({{.FieldName}});
        {{end}}sb.append("}");
        return sb.toString();
    }
{{end}}}
`

// modelLombokTemplate Lombok风格Model模板
const modelLombokTemplate = `package {{.Package}};

import lombok.Data;
{{if .NeedConstructors}}import lombok.NoArgsConstructor;
import lombok.AllArgsConstructor;
{{end}}import java.io.Serializable;
{{range .Imports}}import {{.}};
{{end}}
{{if .UseJsonProperty}}import com.fasterxml.jackson.annotation.JsonProperty;
{{end}}
{{if .TableComment}}/**
 * {{.TableComment}}
 */
{{end}}@Data
{{if .NeedConstructors}}@NoArgsConstructor
@AllArgsConstructor
{{end}}public class {{.ClassName}} implements Serializable {
    private static final long serialVersionUID = 1L;
{{range .Fields}}
{{if .Comment}}    /** {{.Comment}} */
{{end}}{{if $.UseJsonProperty}}{{if $.JsonPropertyUpperCase}}    @JsonProperty("{{title .ColumnName}}")
{{else}}    @JsonProperty("{{.ColumnName}}")
{{end}}{{end}}    private {{.FieldType}} {{.FieldName}};
{{end}}}
`
