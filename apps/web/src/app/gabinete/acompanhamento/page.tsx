import { RoleGate } from "@/features/auth/components/role-gate";
import OfficeDemandBoardPage from "@/features/offices/screens/office-demand-board-page";

export default function Page() {
  return (
    <RoleGate allowed={["councillor", "office_member"]}>
      <OfficeDemandBoardPage />
    </RoleGate>
  );
}
