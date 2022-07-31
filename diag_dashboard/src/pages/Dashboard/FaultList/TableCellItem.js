export default function TableCellItem({row,col,rowSpan,colSpan,children,background}){
    const wrapperStyle={
        gridColumnStart:col,
        gridColumnEnd:col+colSpan,
        gridRowStart:row,
        gridRowEnd:row+rowSpan,
        color:"#FFF",
        zIndex:10,
        fontSize:12,
        background:background,
        opacity:0.8,
        position:'relative',
        border:"1px solid grey",
        padding:5}

    return (
        <div style={wrapperStyle}>
            {children}
        </div>
    )
}