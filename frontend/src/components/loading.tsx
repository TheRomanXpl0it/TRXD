import { RotatingLines } from "react-loader-spinner";

function Loading() {
  return (
    <div className="flex items-center justify-center h-screen bg-background">
      <RotatingLines
        strokeColor="grey"
        width="100"
        visible={true}
      />
    </div>
  );
}

export default Loading;