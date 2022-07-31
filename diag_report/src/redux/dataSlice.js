import { createSlice } from '@reduxjs/toolkit';
// Define the initial state using that type
const initialState = {
    loaded:false,
    //原始数据
    data:{}
}

export const dataSlice = createSlice({
    name: 'data',
    initialState,
    reducers: {
        setData:(state,action) => {
            const {data}=action.payload;
            state.data=data;
            state.loaded=true;
        },
        refreshData:(state,action) => {
            state.loaded=false;
        },
    }
});

// Action creators are generated for each case reducer function
export const { 
    setData,
    refreshData
} = dataSlice.actions

export default dataSlice.reducer