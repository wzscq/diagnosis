import "./globals.css";

export const metadata = {
  title: "vehicle status",
  description: "vehicle status",
};

export default function RootLayout({ children }) {
  
  return (
    <html>
      <body>{children}</body>
    </html>
  );
}
