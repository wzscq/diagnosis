import BottomDrawer from "@/components/BottomDrawer";
import StatusTable from "@/components/StatusTable";
import styles from "./page.module.css";

export default function VehicleStatusLayout({ children }) {
  return (
    <div className={styles.vechileStatusLayout}>
      {children}
      <BottomDrawer>
        <StatusTable />
      </BottomDrawer>
    </div>
  )
}