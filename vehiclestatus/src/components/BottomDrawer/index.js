'use client'
import React, {useState} from 'react';
import {Button,Drawer} from 'antd';
import { UpOutlined } from '@ant-design/icons';
import styles from './page.module.css';

export default function BottomDrawer({children}) {
    const [open, setOpen] = useState(false);

    const onClose = () => {
        setOpen(false);
    }

    const showDrawer = () => {
        setOpen(true);
    }

    return (
        <div className={styles.bottomDrawer}>
            <Button size='small' shape="circle" icon={<UpOutlined />} onClick={showDrawer}/>
            <Drawer
                placement={'bottom'}
                closable={false}
                onClose={onClose}
                open={open}
                key={'bottom'}
                height={430}
            >
                {children}
            </Drawer>
        </div>
    )
}