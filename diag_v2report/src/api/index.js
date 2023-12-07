import { createAsyncThunk } from '@reduxjs/toolkit';
import axios from 'axios';


export const getHost=()=>{
    /*const rootElement=document.getElementById('root');
    const host=rootElement.getAttribute("host");
    console.log("host:"+host);
    return host;*/
    return process.env.SERVER_HOST;
}

const host=getHost()+process.env.REACT_APP_SERVICE_API_PREFIX;

//getReport api
export const getReport = createAsyncThunk(
    'getReport',
    async (params,{ rejectWithValue }) => {
        try{
            const response = await axios({
                url:host+"/reports",
                method:"post",
                params:params
            });
            return response.data
        } catch(err){
            return rejectWithValue(err);
        }
    }
)

export const downloadReport = createAsyncThunk(
    'downloadReport',
    async (params,{ rejectWithValue }) => {
        try{
            console.log('downloadReport',params);
            const response = await axios({
                url:host+"/downloadReport",
                method:"post",
                data:params,
                responseType:'blob'
            });
            return response.data
        } catch(err){
            return rejectWithValue(err);
        }
    }
)

