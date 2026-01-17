package com.example.mapper;

import com.example.model.Photo;
import java.util.List;
import org.apache.ibatis.annotations.Param;

/**
 * PhotoMapper接口
 */
public interface PhotoMapper {
    /**
     * 根据主键删除
     */
    int deleteByPrimaryKey(Integer id);

    /**
     * 插入记录
     */
    int insert(Photo record);

    /**
     * 插入记录（选择性）
     */
    int insertSelective(Photo record);

    /**
     * 根据主键查询
     */
    Photo selectByPrimaryKey(Integer id);

    /**
     * 根据主键更新（选择性）
     */
    int updateByPrimaryKeySelective(Photo record);

    /**
     * 根据主键更新
     */
    int updateByPrimaryKey(Photo record);

}
