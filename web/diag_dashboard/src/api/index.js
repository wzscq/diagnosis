import { createAsyncThunk } from '@reduxjs/toolkit';
import axios from 'axios';


export const getHost=()=>{
    const rootElement=document.getElementById('root');
    const host=rootElement.getAttribute("host");
    console.log("host:"+host);
    return host;
}

const host=getHost()+process.env.REACT_APP_SERVICE_API_PREFIX;

export const getDashboard = createAsyncThunk(
    'getDashboard',
    async (_,{ rejectWithValue }) => {
        try{
            const response = await axios({
                url:host+"/dashboard",
                method:"get",
            });
            return response.data
        } catch(err){
            return rejectWithValue(err);
        }
    }
)