import {Row,Col} from 'antd';
import { useEffect } from 'react';
import { useDispatch } from 'react-redux';

import {setCarInfo} from '../../../redux/reportSlice';

const styleTitle={
    textAlign:"left",
    paddingLeft:20,
    fontSize:18,
    borderBottom:"1px solid #888"
}

const styleLabel={
    textAlign:"right",
    paddingRight:20,
    backgroundColor:"#BBBBBB",
    borderLeft:"1px solid #888",
    borderRight:"1px solid #888",
    borderBottom:"1px solid #888",
    wordBreak:"break-all"
}

const styleContent={
    textAlign:"left",
    paddingLeft:20,
    borderBottom:"1px solid #888",
    borderRight:"1px solid #888",
    wordBreak:"break-all"
}

export default function CarInfo({report,vin}){
    const dispatch=useDispatch()

    const converTime=(time)=>{
        if(time.length<14){
            return time;
        }
        return time.substr(0,4)+'-'+time.substr(4,2)+'-'+time.substr(6,2)+' '+time.substr(8,2)+':'+time.substr(10,2)+':'+time.substr(12,2);
    }

    useEffect(()=>{
        const carInfo=[
            {title:"车辆VIN码",value:vin},
            {title:"项目编号",value:report?.ProjectNum},
            {title:"车辆代码",value:report?.VehicleManagementCode},
            {title:"设备编号",value:report?.DeviceID},
            {title:"文件名",value:report?.DeviceID+"_"+report?.SamplingTime},
            {title:"采集时间",value:report?.SamplingTime},
        ]
        dispatch(setCarInfo(carInfo));
    });

    return (
    <div style={{width:"100%",border:"0px solid #888",borderBottom:0}}>
        <Row>
            <Col span={2} />
            <Col span={20} style={styleTitle}>1、车辆信息</Col>
            <Col span={2} />
        </Row>
        <Row>
            <Col span={2} />
            <Col span={6} style={styleLabel} >车辆VIN码</Col>
            <Col span={14} style={styleContent} >{vin}</Col>
            <Col span={2} />
        </Row>
        <Row>
            <Col span={2} />
            <Col span={6} style={styleLabel} >项目编号</Col>
            <Col span={14} style={styleContent} >{report?.ProjectNum}</Col>
            <Col span={2} />
        </Row>
        <Row>
            <Col span={2} />
            <Col span={6} style={styleLabel} >车辆代码</Col>
            <Col span={14}  style={styleContent} >{report?.VehicleManagementCode}</Col>
            <Col span={2} />
        </Row>
        <Row>
            <Col span={2} />
            <Col span={6} style={styleLabel} >设备编号</Col>
            <Col span={14}  style={styleContent} >{report?.DeviceID}</Col>
            <Col span={2} />
        </Row>
        <Row>
            <Col span={2} />
            <Col span={6} style={styleLabel} >文件名</Col>
            <Col span={14} style={styleContent} >{report?.DeviceID+"_"+report?.SamplingTime}</Col>
            <Col span={2} />
        </Row>
        <Row>
            <Col span={2} />
            <Col span={6} style={styleLabel} >采集时间</Col>
            <Col span={14} style={{...styleContent}} >{converTime(report?.SamplingTime)}</Col>
            <Col span={2} />
        </Row>
    </div>
    )
}