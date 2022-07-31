import {CaretLeftOutlined} from '@ant-design/icons';

export default function GridCellWrapper({row,col,rowSpan,colSpan,children}){
    const wrapperStyle={
        gridColumnStart:col,
        gridColumnEnd:col+colSpan,
        gridRowStart:row,
        gridRowEnd:row+rowSpan,
        color:"#FFF",
        zIndex:10,
        fontSize:12,
        position:'relative',
        border:"1px solid black",
        padding:5}

    return (
        <div style={wrapperStyle}>
            <CaretLeftOutlined style={{color:'#00bcd4',position: 'absolute',left: 0,top: 0,fontSize:'9px',transform:'rotate(45deg)'}} />
            <CaretLeftOutlined style={{color:'#00bcd4',position: 'absolute',right: 0,top: 0,fontSize:'9px',transform:'rotate(135deg)'}} />
            <CaretLeftOutlined style={{color:'#00bcd4',position: 'absolute',right: 0,bottom: 0,fontSize:'9px',transform:'rotate(225deg)'}} />
            <CaretLeftOutlined style={{color:'#00bcd4',position: 'absolute',left: 0,bottom: 0,fontSize:'9px',transform:'rotate(315deg)'}} />
            {children}
        </div>
    )
}