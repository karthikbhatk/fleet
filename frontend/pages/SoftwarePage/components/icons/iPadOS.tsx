import React from "react";

import type { SVGProps } from "react";

const iPadOS = (props: SVGProps<SVGSVGElement>) => {
  // Note: smaller icon on OS table has thicker outline and a smaller Apple logo
  if (props.width === "24") {
    return (
      <svg fill="none" xmlns="http://www.w3.org/2000/svg" {...props}>
        <g transform="translate(4, 4) scale(1.5)">
          <rect
            x="3.875"
            y="2.75"
            width="8.25"
            height="10.5"
            rx="1.25"
            stroke="#515774"
            strokeWidth="1.5"
          />
          <path
            fillRule="evenodd"
            clipRule="evenodd"
            d="M8.832 5.938c.027.22-.079.438-.22.595a.793.793 0 0 1-.613.265c-.03-.213.084-.438.22-.576a.969.969 0 0 1 .613-.285Zm.393 2.324v.004c.075.146.286.367.463.423a.337.337 0 0 1-.01.025l-.008.026c-.049.134-.18.359-.269.477-.176.23-.353.454-.64.466a.834.834 0 0 1-.316-.073c-.101-.039-.205-.079-.37-.077a.902.902 0 0 0-.381.081c-.09.035-.176.07-.299.073-.273.012-.48-.245-.656-.47a3.037 3.037 0 0 1-.27-.473 2.205 2.205 0 0 1-.132-.462v-.008a2.041 2.041 0 0 1-.017-.458v-.008c.013-.146.092-.379.163-.477.176-.277.48-.458.824-.466l.053-.004h.004c.147-.007.29.046.416.092.097.036.183.068.255.066a.915.915 0 0 0 .268-.066c.154-.052.335-.114.512-.1.3.004.59.13.771.367a1.617 1.617 0 0 0-.145.107c0 .004-.005.004-.005.004a.754.754 0 0 0-.282.473v.004a.748.748 0 0 0 .07.454Z"
            fill="#515774"
          />
        </g>
      </svg>
    );
  }
  return (
    <svg fill="none" xmlns="http://www.w3.org/2000/svg" {...props}>
      <g transform="translate(6, 7) scale(0.75)">
        <rect
          x="3.666"
          y=".75"
          width="18"
          height="22.5"
          rx="1.25"
          stroke="#515774"
          strokeWidth="1.5"
        />
        <path
          fillRule="evenodd"
          clipRule="evenodd"
          d="M14.33 7.875c.054.442-.158.876-.44 1.192-.29.324-.767.56-1.225.529-.062-.426.167-.876.44-1.153.3-.315.803-.552 1.226-.568Zm.785 4.65v.007c.15.293.573.735.926.845a.7.7 0 0 1-.018.051.669.669 0 0 0-.017.052c-.097.268-.362.718-.538.955-.353.458-.705.908-1.278.931-.27-.003-.45-.074-.634-.145-.203-.079-.41-.159-.742-.155-.34-.004-.554.08-.76.162-.179.07-.352.139-.597.146-.547.024-.961-.49-1.314-.94-.167-.236-.44-.702-.538-.946-.097-.23-.238-.687-.264-.924v-.016c-.035-.237-.07-.726-.035-.915v-.016c.026-.292.185-.758.326-.955a2.009 2.009 0 0 1 1.648-.932l.106-.008h.009c.294-.013.579.092.83.184.194.071.368.135.51.132.138.003.324-.06.538-.133.306-.104.669-.227 1.023-.199.599.008 1.18.26 1.542.735a3.243 3.243 0 0 0-.29.213c0 .008-.01.008-.01.008-.229.15-.51.552-.564.947v.008c-.053.3.01.67.141.907Z"
          fill="#515774"
        />
      </g>
    </svg>
  );
};

export default iPadOS;
