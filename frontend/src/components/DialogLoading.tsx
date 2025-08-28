import { DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { RotatingLines } from "react-loader-spinner";
function DialogLoading() {
  return (
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Loading...</DialogTitle>
      </DialogHeader>
        <RotatingLines
                strokeColor="grey"
                width="100"
                visible={true}
              />
    </DialogContent>
  );
}

export { DialogLoading };