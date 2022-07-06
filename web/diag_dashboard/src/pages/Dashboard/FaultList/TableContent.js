import TableCellItem from "./TableCellItem"

export default function TableContent({data}){

    const cells=[];
    data.forEach((item,row)=>{
        const background=row%2===0?'#454760':'#666cc2';
        cells.push(
            <TableCellItem col={1} row={row+1} colSpan={1} rowSpan={1} background={background}>
                {item.vehicle_management_code}
            </TableCellItem>
        );
        cells.push(
            <TableCellItem col={2} row={row+1} colSpan={1} rowSpan={1} background={background}>
                {item.project_num}
            </TableCellItem>
        );
        cells.push(
            <TableCellItem col={3} row={row+1} colSpan={1} rowSpan={1} background={background}>
                {item.specifications}
            </TableCellItem>
        );
        cells.push(
            <TableCellItem col={4} row={row+1} colSpan={1} rowSpan={1} background={background}>
                {item.device_number}
            </TableCellItem>
        );
        cells.push(
            <TableCellItem col={5} row={row+1} colSpan={1} rowSpan={1} background={background}>
                {item.time}
            </TableCellItem>
        );
        cells.push(
            <TableCellItem col={6} row={row+1} colSpan={1} rowSpan={1} background={background}>
                {item.type}
            </TableCellItem>
        );
        cells.push(
            <TableCellItem col={7} row={row+1} colSpan={1} rowSpan={1} background={background}>
                {item.status==='0'?(<span style={{color:'red',fontWeight:700}}>Open</span>):(
                    <span style={{color:'#73c0de'}}>Closed</span>
                )}
            </TableCellItem>
        );
        cells.push(
            <TableCellItem col={8} row={row+1} colSpan={1} rowSpan={1} background={background}>
                {item.remark}
            </TableCellItem>
        );
    })

    return (
        <div className="fault-table-content-wrapper">
            <div className="fault-table-content">
                {cells}
            </div>
        </div>
    )
}