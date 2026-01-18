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
{{end}}{{if $.UseJsonProperty}}    @JsonProperty("{{.ColumnName}}")
{{end}}    private {{.FieldType}} {{.FieldName}};
{{end}}
{{range .Fields}}
    public {{.FieldType}} get{{title .FieldName}}() {
        return {{.FieldName}};
    }

    public void set{{title .FieldName}}({{.FieldType}} {{.FieldName}}) {
        this.{{.FieldName}} = {{.FieldName}};
    }
{{end}}}
`

// modelLombokTemplate Lombok风格Model模板
const modelLombokTemplate = `package {{.Package}};

import lombok.Data;
import java.io.Serializable;
{{range .Imports}}import {{.}};
{{end}}
{{if .UseJsonProperty}}import com.fasterxml.jackson.annotation.JsonProperty;
{{end}}
{{if .TableComment}}/**
 * {{.TableComment}}
 */
{{end}}@Data
public class {{.ClassName}} implements Serializable {
    private static final long serialVersionUID = 1L;
{{range .Fields}}
{{if .Comment}}    /** {{.Comment}} */
{{end}}{{if $.UseJsonProperty}}    @JsonProperty("{{.ColumnName}}")
{{end}}    private {{.FieldType}} {{.FieldName}};
{{end}}}
`
