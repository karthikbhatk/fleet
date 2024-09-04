import React from "react";

export const renderEmptySearch = (
  type: "exclamation" | "question" | "check"
) => {
  const renderIconPath = () => {
    switch (type) {
      case "exclamation":
        return (
          <path
            d="m24.24 28.656-1.022-9.774a2.5 2.5 0 1 1 4.973-.003l-1.012 9.775a1.478 1.478 0 0 1-2.94.002h.001Zm-.803 4.951a2.28 2.28 0 1 1 4.562 0 2.28 2.28 0 0 1-4.562 0Z"
            fill="#9C9EDB"
          />
        );
      case "question":
        return (
          <path
            d="M25.836 30.184a1.966 1.966 0 0 1-1.954-2.167c.059-.587.182-1.194.422-1.682.235-.477.995-1.115 2.166-1.979.613-.45 1.04-.87 1.28-1.259.24-.388.36-.87.36-1.446 0-.71-.226-1.268-.677-1.676-.451-.408-.998-.612-1.64-.612-1.274 0-2.13.675-2.566 2.024a1.945 1.945 0 0 1-2.128 1.313 1.932 1.932 0 0 1-1.53-2.614c.416-1.064 1.017-1.933 1.799-2.608 1.29-1.112 2.866-1.668 4.727-1.668 1.86 0 3.491.547 4.777 1.64 1.285 1.094 1.928 2.48 1.928 4.159 0 2.081-1.07 3.822-3.208 5.223-.931.614-1.495 1.093-1.691 1.439-.065.113-.119.259-.162.436a1.955 1.955 0 0 1-1.903 1.478v-.001Zm.317 5.626h-.173a2.172 2.172 0 1 1 0-4.345h.173a2.173 2.173 0 0 1 0 4.345Z"
            fill="#9C9EDB"
          />
        );
      case "check":
        return (
          <path
            d="m31.503 18.45-8.13 10.054a.142.142 0 0 1-.202 0l-3.294-4.383a1.855 1.855 0 1 0-2.641 2.606l4.493 5.58a2.284 2.284 0 0 0 3.275-.026l9.181-11.252c.72-.752.68-1.95-.088-2.652l-.008-.007a1.856 1.856 0 0 0-2.586.08Z"
            fill="#9C9EDB"
          />
        );
      default:
        return null;
    }
  };

  return (
    <svg width="55" height="71" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path
        d="M25.73 49.4c13.123 0 23.762-10.638 23.762-23.761S38.853 1.878 25.73 1.878c-13.123 0-23.761 10.638-23.761 23.76C1.97 38.763 12.608 49.4 25.73 49.4Z"
        fill="#F1F0FF"
        stroke="#fff"
        strokeWidth="1.199"
      />
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M11.669 7.23a23.06 23.06 0 0 1 14.06-4.755c12.792 0 23.162 10.37 23.162 23.162a23.06 23.06 0 0 1-4.758 14.065 23.06 23.06 0 0 1-14.06 4.754c-12.792 0-23.162-10.37-23.162-23.161A23.06 23.06 0 0 1 11.67 7.23Z"
        fill="#fff"
      />
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M11.669 7.809a23.06 23.06 0 0 1 14.06-4.755c12.792 0 23.162 10.37 23.162 23.162a23.06 23.06 0 0 1-4.758 14.065 23.06 23.06 0 0 1-14.06 4.754c-12.792 0-23.162-10.37-23.162-23.161A23.06 23.06 0 0 1 11.67 7.809Z"
        fill="#fff"
      />
      {renderIconPath()}
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M13.687 4.313C1.875 11.132-2.172 26.235 4.647 38.047c6.473 11.211 20.41 15.427 31.91 9.994l4.452 7.712 4.074-2.352-4.48-7.76c10.074-7.34 13.19-21.252 6.818-32.29C40.6 1.542 25.498-2.506 13.687 4.314Zm-6 31.98C1.835 26.158 5.306 13.201 15.44 7.35c10.133-5.85 23.09-2.378 28.941 7.755 5.85 10.133 2.379 23.09-7.755 28.941-10.133 5.85-23.09 2.379-28.94-7.755Z"
        fill="#E3E3FE"
        stroke="#9C9EDB"
        strokeWidth=".882"
      />
      <path
        d="M37.443 52.516a1 1 0 0 1 .366-1.366l4.888-2.822a1 1 0 0 1 1.366.366l9.171 15.886a3 3 0 0 1-1.098 4.098l-1.423.822a3 3 0 0 1-4.098-1.098l-9.172-15.886Z"
        fill="#9C9EDB"
      />
      <path
        d="M37.825 52.296a.559.559 0 0 1 .205-.764l4.887-2.822a.559.559 0 0 1 .764.205L52.853 64.8a2.56 2.56 0 0 1-.937 3.495l-1.423.822a2.559 2.559 0 0 1-3.496-.937l-9.172-15.885Z"
        stroke="#9C9EDB"
        strokeWidth=".882"
      />
    </svg>
  );
};

const EmptySearchQuestion = () => renderEmptySearch("question");

export default EmptySearchQuestion;
