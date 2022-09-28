import {Row,Col} from 'antd';

const styleTitle={
    textAlign:"left",
    paddingLeft:20,
    fontSize:14,
    borderBottom:"1px solid #888",
    borderLeft:"1px solid #888",
    borderRight:"1px solid #888"
}

const styleLabel={
    textAlign:"right",
    paddingRight:20,
    backgroundColor:"#BBBBBB",
    borderLeft:"1px solid #888",
    borderBottom:"1px solid #888",
    wordBreak:"break-all"
}

const styleContent={
    textAlign:"left",
    paddingLeft:20,
    borderRight:"1px solid #888",
    borderBottom:"1px solid #888",
    borderLeft:"1px solid #888",
    wordBreak:"break-all"
}

const styleSubTitle={
    textAlign:"left",
    paddingLeft:20,
    fontSize:18,
    borderBottom:"1px solid #888",
    borderLeft:"1px solid #888",
    borderRight:"1px solid #888",
    backgroundColor:"#BBBBBB"
}

export default function LogisticsInfo({title,logistics,vin}){
    return (
    <div style={{width:"100%",border:"0px solid #888",borderBottom:0, marginTop:0}}>
        <Row >
            <Col span={2} />
            <Col span={20}  style={styleSubTitle}>{title}</Col>
            <Col span={2} />
        </Row>
        <Row>
            <Col span={2} />
            <Col span={20} style={styleTitle} >物流信息</Col>
            <Col span={2} />
        </Row>
        {
            logistics.NameList.map(item=>{
                return (
                    <Row>
                        <Col span={2} />
                        <Col span={6} style={styleLabel} >{item.Name}</Col>
                        <Col span={14} style={styleContent} >{item.value}</Col>
                        <Col span={2} />
                    </Row>
                );
            })
            
        }
        
    </div>
    )
}