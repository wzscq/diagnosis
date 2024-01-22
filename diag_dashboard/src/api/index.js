import { FilterFilled } from '@ant-design/icons';
import { createAsyncThunk } from '@reduxjs/toolkit';
import axios from 'axios';


export const getHost=()=>{
    /*const rootElement=document.getElementById('root');
    const host=rootElement.getAttribute("host");
    console.log("host:"+host);
    return host;*/
    return process.env.REACT_APP_SERVER_HOST;
}

const host=getHost()+process.env.REACT_APP_SERVICE_API_PREFIX;

export const getDashboard = createAsyncThunk(
    'getDashboard',
    async (filters,{ rejectWithValue }) => {
        try{
            const response = await axios({
                url:host+"/dashboard",
                data:filters,
                method:"post",
            });
            return response.data
        } catch(err){
            return rejectWithValue(err);
        }
    }
)