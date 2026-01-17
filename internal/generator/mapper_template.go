package generator

// mapperTemplate Mapper接口模板
const mapperTemplate = `package {{.Package}};

import {{.ModelPackage}}.{{.ModelName}};
import java.util.List;
import org.apache.ibatis.annotations.Param;

/**
 * {{.ModelName}}Mapper接口
 */
public interface {{.MapperName}} {
    /**
     * 根据主键删除
     */
    int deleteByPrimaryKey({{if .PrimaryKey}}{{.PrimaryKey.FieldType}} {{.PrimaryKey.FieldName}}{{else}}Long id{{end}});

    /**
     * 插入记录
     */
    int insert({{.ModelName}} record);

    /**
     * 插入记录（选择性）
     */
    int insertSelective({{.ModelName}} record);

    /**
     * 根据主键查询
     */
    {{.ModelName}} selectByPrimaryKey({{if .PrimaryKey}}{{.PrimaryKey.FieldType}} {{.PrimaryKey.FieldName}}{{else}}Long id{{end}});

    /**
     * 根据主键更新（选择性）
     */
    int updateByPrimaryKeySelective({{.ModelName}} record);

    /**
     * 根据主键更新
     */
    int updateByPrimaryKey({{.ModelName}} record);
{{if .OffsetLimit}}
    /**
     * 分页查询
     */
    List<{{.ModelName}}> selectByPage(@Param("offset") int offset, @Param("limit") int limit);
{{end}}
}
`
