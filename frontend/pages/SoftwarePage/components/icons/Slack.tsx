import React from "react";

import type { SVGProps } from "react";

const Slack = (props: SVGProps<SVGSVGElement>) => (
  <svg xmlns="http://www.w3.org/2000/svg" fill="none" {...props}>
    <path
      fill="#fff"
      stroke="#E2E4EA"
      d="M.5 8A7.5 7.5 0 0 1 8 .5h16A7.5 7.5 0 0 1 31.5 8v16a7.5 7.5 0 0 1-7.5 7.5H8A7.5 7.5 0 0 1 .5 24z"
    />
    <g
      fillRule="evenodd"
      clipPath="url(#Name=slack_clippath)"
      clipRule="evenodd"
    >
      <path
        fill="#36C5F0"
        d="M13.066 5a2.204 2.204 0 0 0 .001 4.408h2.2V7.205A2.204 2.204 0 0 0 13.067 5m0 5.878H7.2A2.202 2.202 0 0 0 5 13.082a2.202 2.202 0 0 0 2.2 2.205h5.866a2.202 2.202 0 0 0 2.2-2.204 2.202 2.202 0 0 0-2.2-2.205"
      />
      <path
        fill="#2EB67D"
        d="M27 13.082a2.202 2.202 0 0 0-2.2-2.204 2.202 2.202 0 0 0-2.2 2.204v2.205h2.2a2.202 2.202 0 0 0 2.2-2.205m-5.867 0V7.204A2.203 2.203 0 0 0 18.933 5a2.202 2.202 0 0 0-2.2 2.204v5.878a2.202 2.202 0 0 0 2.2 2.205 2.202 2.202 0 0 0 2.2-2.205"
      />
      <path
        fill="#ECB22E"
        d="M18.933 27.044a2.202 2.202 0 0 0 2.2-2.204 2.202 2.202 0 0 0-2.2-2.204h-2.2v2.204a2.203 2.203 0 0 0 2.2 2.204m0-5.88H24.8a2.202 2.202 0 0 0 2.2-2.203 2.202 2.202 0 0 0-2.2-2.205h-5.866a2.202 2.202 0 0 0-2.2 2.204 2.202 2.202 0 0 0 2.199 2.205"
      />
      <path
        fill="#E01E5A"
        d="M5 18.96c0 1.217.984 2.204 2.2 2.205a2.202 2.202 0 0 0 2.2-2.204v-2.204H7.2A2.202 2.202 0 0 0 5 18.96m5.867 0v5.88a2.202 2.202 0 0 0 2.2 2.204 2.202 2.202 0 0 0 2.2-2.204v-5.877a2.2 2.2 0 1 0-4.4-.002"
      />
    </g>
    <defs>
      <clipPath id="Name=slack_clippath">
        <path fill="#fff" d="M5 5h22v22.044H5z" />
      </clipPath>
    </defs>
  </svg>
);
export default Slack;
