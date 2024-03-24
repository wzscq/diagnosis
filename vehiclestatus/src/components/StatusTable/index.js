'use client'
import { Table } from 'antd';

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
  

export default function StatusTable() {
    const dataSource=[
        {
            vehicle_code:'V001',
            project_code:'P001',
            test_spec:'TS001',
            device_code:'D001',
            speed_max:53.5,
            speed_avg:23.5,
            mileage:3.5,
            ec_avg:17.5,
            key:'1'
        },
        {
            vehicle_code:'V001',
            project_code:'P001',
            test_spec:'TS001',
            device_code:'D001',
            speed_max:53.5,
            speed_avg:23.5,
            mileage:3.5,
            ec_avg:17.5,
            key:'2'
        },
        {
            vehicle_code:'V001',
            project_code:'P001',
            test_spec:'TS001',
            device_code:'D001',
            speed_max:53.5,
            speed_avg:23.5,
            mileage:3.5,
            ec_avg:17.5,
            key:'3'
        },
        {
            vehicle_code:'V001',
            project_code:'P001',
            test_spec:'TS001',
            device_code:'D001',
            speed_max:53.5,
            speed_avg:23.5,
            mileage:3.5,
            ec_avg:17.5,
            key:'4'
        },
        {
            vehicle_code:'V001',
            project_code:'P001',
            test_spec:'TS001',
            device_code:'D001',
            speed_max:53.5,
            speed_avg:23.5,
            mileage:3.5,
            ec_avg:17.5,
            key:'5'
        },
        {
            vehicle_code:'V001',
            project_code:'P001',
            test_spec:'TS001',
            device_code:'D001',
            speed_max:53.5,
            speed_avg:23.5,
            mileage:3.5,
            ec_avg:17.5,
            key:'6'
        },
        {
            vehicle_code:'V001',
            project_code:'P001',
            test_spec:'TS001',
            device_code:'D001',
            speed_max:53.5,
            speed_avg:23.5,
            mileage:3.5,
            ec_avg:17.5,
            key:'7'
        },
        {
            vehicle_code:'V001',
            project_code:'P001',
            test_spec:'TS001',
            device_code:'D001',
            speed_max:53.5,
            speed_avg:23.5,
            mileage:3.5,
            ec_avg:17.5,
            key:'8'
        }
    ]

    const rowSelection = {
        type:'radio',
        onChange: (selectedRowKeys, selectedRows) => {
            console.log(`selectedRowKeys: ${selectedRowKeys}`, 'selectedRows: ', selectedRows);
        },
        getCheckboxProps: (record) => ({
            disabled: record.name === 'Disabled User',
            // Column configuration not to be checked
            name: record.name,
        }),
    }

    return <Table
        size='small' 
        dataSource={dataSource} 
        columns={columns} 
        pagination={{position:'bottomRight',defaultCurrent:1,total:1}}
        rowSelection={rowSelection}
    />
}