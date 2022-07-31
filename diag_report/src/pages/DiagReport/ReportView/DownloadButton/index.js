import { useDispatch, useSelector } from "react-redux";
import { Button } from "antd";
import { downloadReport } from "../../../../api";

export default function DwonloadButton(){
    const dispatch=useDispatch();
    const report=useSelector(state=>state.report);

    const downReport=()=>{    
        dispatch(downloadReport(report));
    }

    return (
        <Button type='primary' onClick={downReport} key="2">下载报告</Button>
    )
}