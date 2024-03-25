'use client'
import { Table } from 'antd';
import { useRouter } from 'next/navigation'

const columns = [
    {
      title: '车辆编码',
      dataIndex: 'vehicle_code',
      key: 'vehicle_code',
    },
    {
      title: '项目编码',
      dataIndex: 'project_code',
      key: 'project_code',
    },
    {
      title: '试验规范',
      dataIndex: 'test_spec',
      key: 'test_spec',
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
        title: '行驶里程(km)',
        dataIndex: 'mileage',
        key: 'mileage',
    },
    {
        title: '平均能耗(kw)',
        dataIndex: 'ec_avg',
        key: 'ec_avg',
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
        pagination={{position:'bottomRight',defaultCurrent:1,pageSize:8}}
        rowSelection={rowSelection}
    />
}