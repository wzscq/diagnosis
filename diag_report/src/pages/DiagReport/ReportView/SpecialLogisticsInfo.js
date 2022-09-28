import {Row,Col} from 'antd';

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

const styleSubLabel={
    textAlign:"right",
    paddingRight:20,
    fontSize:14,
    backgroundColor:"#BBBBBB",
    borderBottom:"1px solid #888",
    borderLeft:"1px solid #888",
    borderRight:"0px solid #888"
}

export default function SpecialLogisticsInfo({title,logistics,vin}){
    return (
    <>
        {
            Object.keys(logistics).map(key=>{
                return (
                    <Row>
                        <Col span={2} />
                        <Col span={3} style={styleLabel}>
                            {key}
                        </Col>
                        <Col span={3}>
                            <Row>
                                <Col span={24}>
                                    {
                                        Object.keys(logistics[key]).map(itemKey=>{
                                            return (<Row>
                                                <Col span={24} style={styleSubLabel}>{itemKey}</Col>
                                            </Row>);
                                        })
                                    }
                                </Col>
                            </Row>
                        </Col>
                        <Col span={14} >
                            {
                                Object.keys(logistics[key]).map(itemKey=>{
                                    return (<Row >
                                        <Col span={24} style={styleContent}>{logistics[key][itemKey]}</Col>
                                    </Row>);
                                })
                            }
                        </Col>
                        <Col span={2} />
                    </Row>
                );
            })
            
        }
    </>
    );
}