export default function TableCellItem({color,row,col,rowSpan,colSpan,children,background}){
    const wrapperStyle={
        gridColumnStart:col,
        gridColumnEnd:col+colSpan,
        gridRowStart:row,
        gridRowEnd:row+rowSpan,
        color:color,
        zIndex:10,
        fontSize:12,
        background:background,
        opacity:0.8,
        position:'relative',
        border:"1px solid #FAFAFA",
        padding:5}

    return (
        <div style={wrapperStyle}>
            {children}
        </div>
    )
}