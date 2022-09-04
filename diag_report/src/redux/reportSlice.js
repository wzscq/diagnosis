import { createSlice } from '@reduxjs/toolkit';
import { message } from 'antd';

import {downloadReport} from '../api';

// Define the initial state using that type
const initialState = {
    carInfo:[
        {title:"车辆VIN码",value:""},
        {title:"车辆代码",value:""},
        {title:"文件名",value:""},
        {title:"采集时间",value:""},
    ],
    ecuList:[],
    analysisItems:{
        
    },
    fileName:"",
    pending:false
}

const downloadReportFile=(data,fileName)=>{
    //let blob=new Blob([data],{type:`application/octet-stream`});
    var a = document.createElement('a');
    var url = window.URL.createObjectURL(data);
    a.href = url;
    a.download = fileName+".xlsx";
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);
}

export const reportSlice = createSlice({
    name: 'report',
    initialState,
    reducers: {
        setFileName:(state,action)=>{
            state.fileName=action.payload;
        },
        setCarInfo: (state,action) => {
          state.carInfo=action.payload;
        },
        setECUList:(state,action) => {
            state.ecuList=action.payload;
        },
        setAnalysisItem:(state,action)=>{
            const {itmeIndex,item}=action.payload;
            if(state.analysisItems[itmeIndex]===undefined){
                state.analysisItems[itmeIndex]={
                    item:{},
                    signalChart:{}
                }
            }
            state.analysisItems[itmeIndex].item=item
        },
        setAnalysisChart:(state,action)=>{
            console.log('setAnalysisChart',action.payload);
            const {itmeIndex,chartIndex,chart}=action.payload;
            if(state.analysisItems[itmeIndex]===undefined){
                state.analysisItems[itmeIndex]={
                    item:{},
                    signalChart:{}
                }
            }
            state.analysisItems[itmeIndex].signalChart[chartIndex]=chart
        }
    },
    extraReducers: (builder) => {
        builder.addCase(downloadReport.pending, (state, action) => {
            state.pending=true;
        });
        builder.addCase(downloadReport.fulfilled, (state, action) => {
            console.log('getReport fulfilled',action);
            downloadReportFile(action.payload,state.fileName);
            state.pending=false;
        });
        builder.addCase(downloadReport.rejected , (state, action) => {
            state.pending=false;
            message.error("获取报告出错！");
        });
    }
});

// Action creators are generated for each case reducer function
export const { 
    setCarInfo,
    setECUList,
    setAnalysisItem,
    setAnalysisChart,
    setFileName
} = reportSlice.actions

export default reportSlice.reducer