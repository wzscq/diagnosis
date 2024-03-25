import BottomDrawer from "@/components/BottomDrawer";
import StatusTable from "@/components/StatusTable";
import {getVechicleStatus} from "@/api";

import styles from "./page.module.css";

export default async function VehicleStatusLayout({ children }) {
  const data = await getVechicleStatus();

  return (
    <div className={styles.vechileStatusLayout}>
      {children}
      <BottomDrawer>
        <StatusTable dataSource={data?.result?.list??[]}/>
      </BottomDrawer>
    </div>
  )
}