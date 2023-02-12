import {Row,Col} from 'antd';

import AnalysisItem from "./AnalysisItem";

const styleTitle={
    textAlign:"left",
    paddingLeft:20,
    fontSize:18,
    borderBottom:"0px solid #888"
}

export default function Analysis({report,vin,itemIndex}){
    const {Records}=report;
    const items=Records.map((rec,index)=>{
        return (<AnalysisItem itmeIndex={itemIndex+''} report={report} rec={rec} vin={vin} />);
    });

    return (
        <div style={{width:"100%",border:"0px solid #EEEEEE",borderBottom:0}}>
            <Row>
                <Col span={2} />
                <Col span={20} style={styleTitle}>DTC({report.ReportID})</Col>
                <Col span={2} />
            </Row>
            {
                items
            }
        </div>
    );
}