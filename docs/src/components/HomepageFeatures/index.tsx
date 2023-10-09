import React from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  Svg: React.ComponentType<React.ComponentProps<'svg'>>;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Manage Tech Resources',
    Svg: require('@site/static/img/undraw_docusaurus_mountain.svg').default,
    description: (
      <>
        terediX fetch, monitor your tech resources from various sources and make it available for you in one place.
      </>
    ),
  },
  {
    title: 'Solve Real World Problems',
    Svg: require('@site/static/img/undraw_docusaurus_mountain.svg').default,
    description: (
      <>
        terediX helps you to solve real world problems like, resource visibility, resource management, resource monitoring, resource reporting etc.
      </>
    ),
  },
  {
    title: 'Resource Visibility',
    Svg: require('@site/static/img/undraw_docusaurus_tree.svg').default,
    description: (
      <>
        Explore your tech resources in one place. Greater visibility of your resources will help you to make better decisions.
      </>
    ),
  },
  {
    title: 'Easy to Configure',
    Svg: require('@site/static/img/undraw_docusaurus_react.svg').default,
    description: (
      <>
        Configure teredix with a simple configuration file. Run the tools, it will start working it's job.
      </>
    ),
  },
  {
    title: 'Multiple Deployment Options',
    Svg: require('@site/static/img/undraw_docusaurus_react.svg').default,
    description: (
      <>
        Deploy terediX in your own machine or deploy in Kubernetes or just run it as a docker container.
      </>
    ),
  },
  {
    title: 'Monitoring & Reporting',
    Svg: require('@site/static/img/undraw_docusaurus_react.svg').default,
    description: (
      <>
        TerediX exposes prometheus metrics and it already provides some built-in Grafana dashboard to monitor your resources.
      </>
    ),
  },
  {
    title: 'Fully Open Source',
    Svg: require('@site/static/img/undraw_docusaurus_react.svg').default,
    description: (
      <>
        terediX is fully open source and community driven project. You can contribute to the project and make it better.
      </>
    ),
  },
  {
    title: 'Support Major Resource Types',
    Svg: require('@site/static/img/undraw_docusaurus_tree.svg').default,
    description: (
        <>
          terediX can connect to major resource providers such as AWS, GitHub
        </>
    ),
  },
];

function Feature({title, Svg, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <Svg className={styles.featureSvg} role="img" />
      </div>
      <div className="text--center padding-horiz--md">
        <h3>{title}</h3>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): JSX.Element {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
