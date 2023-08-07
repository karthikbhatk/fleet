import React from "react";
import { COLORS, Colors } from "styles/var/colors";
import { ICON_SIZES, IconSizes } from "styles/var/icon_sizes";

interface ILinuxCircledProps {
  size?: IconSizes;
  iconColor?: Colors;
  bgColor?: Colors;
}

const LinuxCircled = ({
  size = "extra-large",
  iconColor = "ui-fleet-black-75", // default grey
  bgColor = "ui-blue-10", // default light blue
}: ILinuxCircledProps) => {
  return (
    <svg
      width={ICON_SIZES[size]}
      height={ICON_SIZES[size]}
      viewBox="0 0 48 48"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <circle cx="24" cy="24" r="24" fill={COLORS[bgColor]} />
      <rect
        width="24"
        height="24"
        transform="translate(12.5 12)"
        fill="#F1F0FF"
      />
      <path
        d="M36.4875 28.44C35.4131 24.19 33.1513 21.99 31.8847 21.05C31.8169 20.96 31.749 20.87 31.6812 20.78C31.8847 20.16 31.9978 19.5 31.9978 18.82C31.9978 15.05 28.6391 12 24.4887 12C20.3384 12 16.9683 15.06 16.9683 18.82C16.9683 19.5 17.0814 20.16 17.285 20.78C17.1945 20.89 17.1153 21 17.0362 21.11C15.7583 22.08 13.5757 24.28 12.5126 28.44C12.4674 28.62 12.5465 28.82 12.7049 28.94C12.8745 29.06 13.1007 29.08 13.2929 29.01C13.881 28.77 14.6613 28.41 15.3851 27.91C15.7696 30.22 17.0362 32.23 18.8456 33.6H17.4885C16.7421 33.6 16.1315 34.14 16.1315 34.8C16.1315 35.46 16.7421 36 17.4885 36H31.4211C32.1675 36 32.7781 35.46 32.7781 34.8C32.7781 34.14 32.1675 33.6 31.4211 33.6H30.1092C31.9187 32.23 33.1966 30.2 33.5698 27.88C34.3048 28.39 35.0965 28.76 35.6958 29C35.8881 29.08 36.1143 29.05 36.2839 28.93C36.4535 28.82 36.5327 28.63 36.4875 28.44ZM24.4774 32.73C20.9377 32.73 18.0653 29.91 18.0653 26.44C18.0653 24.66 18.823 22.41 20.033 20.65C19.3997 20.16 19.0039 19.44 19.0039 18.65C19.0039 17.17 20.361 15.98 22.0234 15.98C23.0299 15.98 23.912 16.41 24.4661 17.08C25.0202 16.41 25.9023 15.98 26.9088 15.98C28.5825 15.98 29.9283 17.18 29.9283 18.65C29.9283 19.45 29.5325 20.16 28.8992 20.65C30.1092 22.41 30.8669 24.65 30.8669 26.44C30.8896 29.92 28.0171 32.73 24.4774 32.73Z"
        fill={COLORS[iconColor]}
      />
      <path
        d="M26.2415 19.78L24.5565 19.32C24.5113 19.31 24.4547 19.31 24.3982 19.32L22.7132 19.78C22.634 19.8 22.5662 19.86 22.5435 19.93C22.5209 20 22.5322 20.08 22.5775 20.14L24.2625 22.15C24.3077 22.21 24.3869 22.24 24.466 22.24C24.5452 22.24 24.6244 22.21 24.6696 22.15L26.3546 20.14C26.3999 20.08 26.4225 20 26.3886 19.93C26.3886 19.86 26.3207 19.81 26.2415 19.78Z"
        fill={COLORS[iconColor]}
      />
      <path
        d="M23.4938 18.6898C23.4938 19.0398 23.1771 19.3198 22.7813 19.3198C22.3855 19.3198 22.0688 19.0398 22.0688 18.6898C22.0688 18.3398 22.3855 18.0598 22.7813 18.0598C23.1771 18.0498 23.4938 18.3398 23.4938 18.6898Z"
        fill={COLORS[iconColor]}
      />
      <path
        d="M26.8976 18.6898C26.8976 19.0398 26.5809 19.3198 26.1851 19.3198C25.7893 19.3198 25.4727 19.0398 25.4727 18.6898C25.4727 18.3398 25.7893 18.0598 26.1851 18.0598C26.5696 18.0498 26.8976 18.3398 26.8976 18.6898Z"
        fill={COLORS[iconColor]}
      />
    </svg>
  );
};

export default LinuxCircled;
