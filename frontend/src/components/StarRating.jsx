export default function StarRating({ value = 0, size = "" }) {
  return (
    <span className={`stars ${size}`}>
      {[1, 2, 3, 4, 5].map((i) => (
        <span key={i} className={i <= value ? "" : "off"}>★</span>
      ))}
    </span>
  );
}
