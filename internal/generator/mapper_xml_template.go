package generator

// mapperXMLTemplate Mapper XML模板
const mapperXMLTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" 
"http://mybatis.org/dtd/mybatis-3-mapper.dtd">
<mapper namespace="{{.Namespace}}">
    <!-- ResultMap -->
    <resultMap id="BaseResultMap" type="{{.ModelType}}">
{{if .PrimaryKey}}        <id column="{{.PrimaryKey.ColumnName}}" jdbcType="{{.PrimaryKey.JdbcType}}" property="{{.PrimaryKey.FieldName}}" />
{{end}}{{range .Columns}}{{if ne .ColumnName $.PrimaryKey.ColumnName}}        <result column="{{.ColumnName}}" jdbcType="{{.JdbcType}}" property="{{.FieldName}}" />
{{end}}{{end}}    </resultMap>

    <!-- 基础列 -->
    <sql id="Base_Column_List">
        {{range $index, $col := .Columns}}{{if $index}}, {{end}}{{$col.ColumnName}}{{end}}
    </sql>

    <!-- 根据主键查询 -->
    <select id="selectByPrimaryKey" parameterType="{{.PrimaryKey.JavaType}}" resultMap="BaseResultMap">
        SELECT <include refid="Base_Column_List" />
        FROM {{.TableName}}
        WHERE {{.PrimaryKey.ColumnName}} = #{{"{"}}{{.PrimaryKey.FieldName}},jdbcType={{.PrimaryKey.JdbcType}}{{"}"}}
    </select>

    <!-- 插入 -->
    <insert id="insert" parameterType="{{.ModelType}}"{{if .UseGeneratedKeys}} useGeneratedKeys="true" keyProperty="{{.GenerateKeys}}"{{end}}>
        INSERT INTO {{.TableName}} (
            {{range $index, $col := .Columns}}{{if ne $col.ColumnName $.PrimaryKey.ColumnName}}{{if $index}}, {{end}}{{$col.ColumnName}}{{end}}{{end}}
        )
        VALUES (
            {{range $index, $col := .Columns}}{{if ne $col.ColumnName $.PrimaryKey.ColumnName}}{{if $index}}, {{end}}#{{"{"}}{{$col.FieldName}},jdbcType={{$col.JdbcType}}{{"}"}}{{end}}{{end}}
        )
    </insert>

    <!-- 选择性插入 -->
    <insert id="insertSelective" parameterType="{{.ModelType}}"{{if .UseGeneratedKeys}} useGeneratedKeys="true" keyProperty="{{.GenerateKeys}}"{{end}}>
        INSERT INTO {{.TableName}}
        <trim prefix="(" suffix=")" suffixOverrides=",">
{{range .Columns}}{{if ne .ColumnName $.PrimaryKey.ColumnName}}            <if test="{{.FieldName}} != null">
                {{.ColumnName}},
            </if>
{{end}}{{end}}        </trim>
        <trim prefix="values (" suffix=")" suffixOverrides=",">
{{range .Columns}}{{if ne .ColumnName $.PrimaryKey.ColumnName}}            <if test="{{.FieldName}} != null">
                #{{"{"}}{{.FieldName}},jdbcType={{.JdbcType}}{{"}"}},
            </if>
{{end}}{{end}}        </trim>
    </insert>

    <!-- 根据主键更新 -->
    <update id="updateByPrimaryKey" parameterType="{{.ModelType}}">
        UPDATE {{.TableName}}
        SET {{range $index, $col := .Columns}}{{if ne $col.ColumnName $.PrimaryKey.ColumnName}}{{if $index}},
            {{end}}{{$col.ColumnName}} = #{{"{"}}{{$col.FieldName}},jdbcType={{$col.JdbcType}}{{"}"}}{{end}}{{end}}
        WHERE {{.PrimaryKey.ColumnName}} = #{{"{"}}{{.PrimaryKey.FieldName}},jdbcType={{.PrimaryKey.JdbcType}}{{"}"}}
    </update>

    <!-- 选择性更新 -->
    <update id="updateByPrimaryKeySelective" parameterType="{{.ModelType}}">
        UPDATE {{.TableName}}
        <set>
{{range .Columns}}{{if ne .ColumnName $.PrimaryKey.ColumnName}}            <if test="{{.FieldName}} != null">
                {{.ColumnName}} = #{{"{"}}{{.FieldName}},jdbcType={{.JdbcType}}{{"}"}},
            </if>
{{end}}{{end}}        </set>
        WHERE {{.PrimaryKey.ColumnName}} = #{{"{"}}{{.PrimaryKey.FieldName}},jdbcType={{.PrimaryKey.JdbcType}}{{"}"}}
    </update>

    <!-- 根据主键删除 -->
    <delete id="deleteByPrimaryKey" parameterType="{{.PrimaryKey.JavaType}}">
        DELETE FROM {{.TableName}}
        WHERE {{.PrimaryKey.ColumnName}} = #{{"{"}}{{.PrimaryKey.FieldName}},jdbcType={{.PrimaryKey.JdbcType}}{{"}"}}
    </delete>
{{if .OffsetLimit}}
    <!-- 分页查询 -->
    <select id="selectByPage" resultMap="BaseResultMap">
        SELECT <include refid="Base_Column_List" />
        FROM {{.TableName}}
        LIMIT #{offset}, #{limit}
    </select>
{{end}}
{{if .UseBatchInsert}}
    <!-- 批量插入 -->
    <insert id="insertBatch" parameterType="java.util.List"{{if .UseGeneratedKeys}} useGeneratedKeys="true" keyProperty="{{.GenerateKeys}}"{{end}}>
        INSERT INTO {{.TableName}} (
            {{range $index, $col := .Columns}}{{if ne $col.ColumnName $.PrimaryKey.ColumnName}}{{if $index}}, {{end}}{{$col.ColumnName}}{{end}}{{end}}
        )
        VALUES
        <foreach collection="list" item="item" separator=",">
            ({{range $index, $col := .Columns}}{{if ne $col.ColumnName $.PrimaryKey.ColumnName}}{{if $index}}, {{end}}#{{"{"}}{{"item."}}{{$col.FieldName}},jdbcType={{$col.JdbcType}}{{"}"}}{{end}}{{end}})
        </foreach>
    </insert>
{{end}}
{{if .UseBatchUpdate}}
    <!-- 批量更新 -->
    <update id="updateBatch" parameterType="java.util.List">
        <foreach collection="list" item="item" separator=";">
            UPDATE {{.TableName}}
            SET {{range $index, $col := .Columns}}{{if ne $col.ColumnName $.PrimaryKey.ColumnName}}{{if $index}},
                {{end}}{{$col.ColumnName}} = #{{"{"}}{{"item."}}{{$col.FieldName}},jdbcType={{$col.JdbcType}}{{"}"}}{{end}}{{end}}
            WHERE {{.PrimaryKey.ColumnName}} = #{{"{"}}{{"item."}}{{.PrimaryKey.FieldName}},jdbcType={{.PrimaryKey.JdbcType}}{{"}"}}
        </foreach>
    </update>
{{end}}
</mapper>
`
