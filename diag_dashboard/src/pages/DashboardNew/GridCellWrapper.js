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
        border:"5px solid #F0F0F0",
        padding:5}

    return (
        <div style={wrapperStyle}>
            {children}
        </div>
    )
}