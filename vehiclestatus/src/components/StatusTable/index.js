'use client'
import { Table } from 'antd';
import { useRouter } from 'next/navigation'

const columns = [
    {
      title: '创建时间',
      dataIndex: 'create_time',
      key: 'create_time',
    },
    {
        title: '设备编号',
        dataIndex: 'device_code',
        key: 'device_code',
    },
    {
        title: '最大车速(km/h)',
        dataIndex: 'speed_max',
        key: 'speed_max',
    },
    {
        title: '平均车速(km/h)',
        dataIndex: 'speed_avg',
        key: 'speed_avg',
    },
    {
        title: '平均能耗(kwh/100km)',
        dataIndex: 'ec_avg',
        key: 'ec_avg',
    },
    {
        title: 'NOP激活行驶里程(km)',
        dataIndex: 'nop_mileage',
        key: 'nop_mileage',
    },
    {
        title: 'NOP退出次数',
        dataIndex: 'nop_disabledTimes',
        key: 'nop_disabledTimes',
    },
    {
        title: 'TJA激活行驶里程(km)',
        dataIndex: 'tja_mileage',
        key: 'tja_mileage',
    },
    {
        title: 'ACC激活行驶里程(km)',
        dataIndex: 'acc_mileage',
        key: 'acc_mileage',
    }
];
  

export default function StatusTable({dataSource}) {
    const router = useRouter()
    const rowSelection = {
        type:'radio',
        onChange: (selectedRowKeys, selectedRows) => {
            console.log(`selectedRowKeys: ${selectedRowKeys}`, 'selectedRows: ', selectedRows);
            if(selectedRowKeys.length>0){
                router.push('/vehiclestatus/'+selectedRowKeys[0])
            }
        },
        getCheckboxProps: (record) => ({
            disabled: record.name === 'Disabled User',
            // Column configuration not to be checked
            name: record.name,
        }),
    }

    dataSource.forEach(element => {
        element['key'] = element['id']
    });

    return <Table
        size='small' 
        dataSource={dataSource} 
        columns={columns} 
        pagination={{position:'bottomRight',defaultCurrent:1,pageSize:8,showSizeChanger:false}}
        rowSelection={rowSelection}
    />
}