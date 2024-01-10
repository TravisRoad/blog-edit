"use client";

import React from "react";
import { serialize } from "next-mdx-remote/serialize";
import { MDXRemote, MDXRemoteSerializeResult } from "next-mdx-remote";
// import { MDXRemote } from "next-mdx-remote/rsc";

const Foo = () => {
  return (
    <>
      <div className="text-red-400"> this is a jsx component </div>
    </>
  );
};

const Mdx: React.FC<{ content?: string }> = ({ content }) => {
  const [result, setResult] = React.useState<MDXRemoteSerializeResult>();
  React.useEffect(() => {
    console.log(content);
    if (!content) {
      return;
    }

    serialize(content, {
      mdxOptions: {
        // https://github.com/hashicorp/next-mdx-remote/issues/350#issuecomment-1803930061
        development: process.env.NODE_ENV === "development",
      },
      parseFrontmatter: true,
    }).then((data: MDXRemoteSerializeResult) => {
      setResult(data);
    });
  }, [content]);
  if (result === undefined) {
    return <div>idk</div>;
  }

  return (
    <div className="prose prose-stone">
      <h1>{result.frontmatter.title as string}</h1>
      <ul>
        {(result.frontmatter.list as string[]).map((i) => (
          <li>{i}</li>
        ))}
      </ul>
      <MDXRemote {...result} components={{ Foo }} lazy={true} />
    </div>
  );

  // const [component, setComponent] =
  //   React.useState<
  //     React.ReactElement<any, string | React.JSXElementConstructor<any>>
  //   >();

  // React.useEffect(() => {
  //    <MDXRemote
  //     source={content}
  //     components={{ Foo }}
  //     options={{
  //       mdxOptions: { rehypePlugins: [], remarkPlugins: [] },
  //       scope: {},
  //     }}
  //   />.then((data) => setComponent(data));
  // }, []);
};

export default Mdx;
