package com.example.model;

import java.io.Serializable;


/**
 * 
 */
public class Photo implements Serializable {
    private static final long serialVersionUID = 1L;

    /**  */
    private Integer id;

    /**  */
    private String exposuretime;

    /**  */
    private String iso;

    /**  */
    private String fnumber;

    /**  */
    private String focallength;

    /**  */
    private String model;

    /**  */
    private String origindate;

    /**  */
    private String file;


    public Integer getId() {
        return id;
    }

    public void setId(Integer id) {
        this.id = id;
    }

    public String getExposuretime() {
        return exposuretime;
    }

    public void setExposuretime(String exposuretime) {
        this.exposuretime = exposuretime;
    }

    public String getIso() {
        return iso;
    }

    public void setIso(String iso) {
        this.iso = iso;
    }

    public String getFnumber() {
        return fnumber;
    }

    public void setFnumber(String fnumber) {
        this.fnumber = fnumber;
    }

    public String getFocallength() {
        return focallength;
    }

    public void setFocallength(String focallength) {
        this.focallength = focallength;
    }

    public String getModel() {
        return model;
    }

    public void setModel(String model) {
        this.model = model;
    }

    public String getOrigindate() {
        return origindate;
    }

    public void setOrigindate(String origindate) {
        this.origindate = origindate;
    }

    public String getFile() {
        return file;
    }

    public void setFile(String file) {
        this.file = file;
    }

}
